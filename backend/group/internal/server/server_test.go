package server_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"testing"

	grpb "github.com/Egot3/supel/backend/contracts/group"
	"github.com/egot3/supel/backend/group/internal/database/repositories/curator"
	"github.com/egot3/supel/backend/group/internal/database/repositories/group"
	"github.com/egot3/supel/backend/group/internal/database/repositories/member"
	"github.com/egot3/supel/backend/group/internal/interceptors"
	"github.com/egot3/supel/backend/group/internal/logctx"
	"github.com/egot3/supel/backend/group/internal/server"
	testutils "github.com/egot3/supel/backend/group/internal/testUtils"
	"github.com/egot3/supel/backend/group/internal/types"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func startProdServer(t *testing.T, i do.Injector) *grpc.ClientConn {
	t.Helper()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer(grpc.UnaryInterceptor(interceptors.LoggingUnaryInterceptor(logger)))

	logctx.WithLogger(t.Context(), logger)

	server, err := server.NewGroupService(i)
	require.NoError(t, err)

	grpb.RegisterGroupServiceServer(srv, server)
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Logf("bufconn server stopped: %v", err)
		}
	}()

	t.Cleanup(func() { srv.GracefulStop(); lis.Close() })

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return conn
}

func TestGroup_Group_CreationAndDeletion(t *testing.T) {
	i := testutils.NewTestInjector(t)

	do.Provide(i, group.NewGroupRepository)
	do.Provide(i, curator.NewGroupRepository)
	do.Provide(i, member.NewMemberRepository)
	do.Provide(i, testutils.AllowAllStub)
	do.Provide(i, testutils.AllowAllUsers)

	conn := startProdServer(t, i)
	client := grpb.NewGroupServiceClient(conn)

	grRepo := do.MustInvoke[group.GroupRepository](i)

	testCases := []struct {
		desc        string
		errExpected bool
		name        string
		description *string
	}{
		{
			desc:        "create valid group",
			name:        "test_name",
			description: nil,
			errExpected: false,
		},
		{
			desc:        "create group with empty name",
			name:        "",
			errExpected: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := client.CreateGroup(context.Background(), &grpb.CreateGroupRequest{
				Name:        tC.name,
				Description: tC.description,
				CuratorUUID: uuid.Nil.String(),
				GroupType:   grpb.GroupType_CLUB,
			})
			if tC.errExpected {
				t.Log(err.Error())
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			var groupUUIDString string
			require.NoError(t, func() error {
				resp, err := client.ListGroups(context.Background(), &grpb.ListGroupsRequest{
					Page:      0,
					Size:      2,
					GroupType: grpb.GroupType_UNSPECIFIED,
				})
				if err != nil {
					return err
				}

				groups, total, _ := grRepo.ListGroups(context.Background(), 0, 1, types.UNSPECIFIED, bun.OrderAsc)
				t.Logf("groups: %v\ntotal: %d", groups, total)
				if resp.Total != 1 {
					return fmt.Errorf("%d groups created from 1", resp.Total)
				}

				assert.Equal(t, tC.name, resp.Groups[0].Name)
				groupUUIDString = resp.Groups[0].UUID
				return nil
			}())

			_, err = client.DeleteGroup(context.Background(), &grpb.DeleteGroupRequest{
				GroupUUID: groupUUIDString,
			})
			require.NoError(t, func() error {
				resp, err := client.ListGroups(context.Background(), &grpb.ListGroupsRequest{
					Page: 0,
					Size: 1,
				})
				if err != nil {
					return err
				}
				assert.Equal(t, 0, int(resp.Total))
				return nil
			}())
		})
	}
}

func TestGroup_Curator_To_GroupAssignmentAndRevokation(t *testing.T) {
	i := testutils.NewTestInjector(t)

	do.Provide(i, group.NewGroupRepository)
	do.Provide(i, curator.NewGroupRepository)
	do.Provide(i, member.NewMemberRepository)
	do.Provide(i, testutils.AllowAllStub)
	do.Provide(i, testutils.AllowAllUsers)

	testCurUUID := uuid.New()
	/* testSenUUID := uuid.New() */
	var testGrUUID uuid.UUID

	curRepo := do.MustInvoke[curator.CuratorRepository](i)
	grRepo := do.MustInvoke[group.GroupRepository](i)

	require.NoError(t, grRepo.CreateGroup(context.Background(), uuid.Nil, t.Name(), nil, types.CLUB))
	require.NoError(t, func() error {
		groups, total, err := grRepo.ListGroups(context.Background(), 0, 2, types.CLUB, bun.OrderDesc)
		if err != nil {
			return err
		}
		if total != 1 {
			return fmt.Errorf("%d groups created from 1", total)
		}

		t.Logf("got %d groups", len(groups))
		t.Logf("groups: %v", groups)
		testGrUUID = groups[0].UUID
		t.Logf("new Test gr UUID: %v", testGrUUID)
		return nil
	}())
	require.NoError(t, curRepo.AddCurator(context.Background(), uuid.Nil, testCurUUID, testGrUUID))

	conn := startProdServer(t, i)
	client := grpb.NewGroupServiceClient(conn)

	t.Run("known curator and group", func(t *testing.T) {
		t.Logf("Test gr uuid: %v", testGrUUID.String())
		resp, err := client.GroupsCurators(context.Background(), &grpb.GroupsCuratorsRequest{
			GroupUUID: testGrUUID.String(),
		})
		require.NoError(t, err)
		t.Logf("get curs resp: %v", resp.Curators)
		protoUUID, err := uuid.Parse(resp.GetCurators()[0].Uuid)
		require.NoError(t, err)
		assert.Equal(t, testCurUUID, protoUUID)

		_, err = client.RevokeCuratorFromGroup(t.Context(), &grpb.RevokeCuratorFromGroupRequest{
			CuratorUUID: testCurUUID.String(),
			GroupUUID:   testGrUUID.String(),
		})
	})

	testCases := []struct {
		desc        string
		curatorUUID uuid.UUID
		groupUUID   uuid.UUID
	}{
		{
			desc:        "unknown group known curator",
			curatorUUID: testCurUUID,
			groupUUID:   uuid.New(),
		},
		{
			desc:        "known group unknown curator",
			curatorUUID: uuid.New(),
			groupUUID:   testGrUUID,
		},
		{
			desc:        "unknown group unknown curator",
			curatorUUID: uuid.New(),
			groupUUID:   uuid.New(),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			_, err := client.AssignCuratorToGroup(t.Context(), &grpb.AssignCuratorToGroupRequest{
				CuratorUUID: tC.curatorUUID.String(),
				GroupUUID:   tC.groupUUID.String(),
			})

			require.Error(t, err)
		})
	}
}

/* guess no CTE for sqlite
func TestGroup_Curator_To_SeniorAssignment(t *testing.T)
*/

func TestGroup_Member_To_GroupAssignmentAndRevokation(t *testing.T) {
	i := testutils.NewTestInjector(t)

	do.Provide(i, group.NewGroupRepository)
	do.Provide(i, curator.NewGroupRepository)
	do.Provide(i, member.NewMemberRepository)
	do.Provide(i, testutils.AllowAllStub)
	do.Provide(i, testutils.AllowAllUsers)

	var testGrUUID uuid.UUID

	grRepo := do.MustInvoke[group.GroupRepository](i)

	require.NoError(t, grRepo.CreateGroup(context.Background(), uuid.Nil, t.Name(), nil, types.CLUB))
	require.NoError(t, func() error {
		groups, total, err := grRepo.ListGroups(context.Background(), 0, 2, types.CLUB, bun.OrderDesc)
		if err != nil {
			return err
		}
		if total != 1 {
			return fmt.Errorf("%d groups created from 1", total)
		}

		t.Logf("got %d groups", len(groups))
		t.Logf("groups: %v", groups)
		testGrUUID = groups[0].UUID
		t.Logf("new Test gr UUID: %v", testGrUUID)
		return nil
	}())

	conn := startProdServer(t, i)
	client := grpb.NewGroupServiceClient(conn)

	testCases := []struct {
		desc        string
		groupUUID   uuid.UUID
		errExpected bool
	}{
		{
			desc:        "add and remove new member to-from known group",
			groupUUID:   testGrUUID,
			errExpected: false,
		},
		{
			desc:        "add and remove new member to-from unknown group",
			groupUUID:   uuid.Nil,
			errExpected: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			testMemberUUID := uuid.New()
			_, err := client.AddMember(context.Background(), &grpb.AddMemberRequest{
				MemberUUID: testMemberUUID.String(),
				GroupUUID:  tC.groupUUID.String(),
			})
			t.Logf("error: %v", err)
			if tC.errExpected {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			resp, err := client.ListMembers(context.Background(), &grpb.ListMembersRequest{
				Page:      0,
				Size:      1,
				GroupUUID: tC.groupUUID.String(),
			})
			require.NoError(t, err)
			assert.Equal(t, 1, int(resp.Total))
			assert.Equal(t, testMemberUUID.String(), resp.Members[0].Uuid)

			_, err = client.RemoveMember(context.Background(), &grpb.RemoveMemberRequest{
				MemberUUID: testMemberUUID.String(),
				GroupUUID:  tC.groupUUID.String(),
			})
			require.NoError(t, err)
		})
	}
}
