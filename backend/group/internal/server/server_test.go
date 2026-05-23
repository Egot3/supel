package server_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	grpb "github.com/Egot3/supel/backend/contracts/group"
	"github.com/egot3/supel/backend/group/internal/database/repositories/curator"
	"github.com/egot3/supel/backend/group/internal/database/repositories/group"
	"github.com/egot3/supel/backend/group/internal/database/repositories/member"
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

	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer()

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

func TestGroup_CreateCurator(t *testing.T) {
	i := testutils.NewTestInjector(t)

	do.Provide(i, group.NewGroupRepository)
	do.Provide(i, curator.NewGroupRepository)
	do.Provide(i, member.NewMemberRepository)
	do.Provide(i, testutils.AllowAllStub)

	testCurUUID := uuid.New()
	/* testSenUUID := uuid.New() */
	var testGrUUID uuid.UUID

	curRepo := do.MustInvoke[curator.CuratorRepository](i)
	grRepo := do.MustInvoke[group.GroupRepository](i)

	require.NoError(t, grRepo.CreateGroup(context.Background(), uuid.Nil, t.Name(), nil, types.CLUB))
	require.NoError(t, func() error {
		groups, total, err := grRepo.ListGroups(context.Background(), 0, 1, types.CLUB, bun.OrderDesc)
		if err != nil {
			return err
		}
		if total != 1 {
			return fmt.Errorf("%d groups created from 1", total)
		}
		testGrUUID = groups[1].UUID
		return nil
	}())
	require.NoError(t, curRepo.AddCurator(context.Background(), uuid.Nil, testCurUUID, testGrUUID))
	require.NoError(t, curRepo.AssignCuratorToGroup(context.Background(), testCurUUID, testGrUUID))

	conn := startProdServer(t, i)
	client := grpb.NewGroupServiceClient(conn)

	t.Run("known curator and group", func(t *testing.T) {
		t.Logf("Test gr uuid: %v", testGrUUID.String())
		resp, err := client.GroupsCurators(context.Background(), &grpb.GroupsCuratorsRequest{
			GroupUUID: testGrUUID.String(),
		})
		require.NoError(t, err)
		t.Logf("get curs resp: %v", resp.Curators)
		assert.Equal(t, testCurUUID, resp.GetCurators()[0].Uuid)
	})
}
