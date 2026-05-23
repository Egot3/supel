package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	grpb "github.com/Egot3/supel/backend/contracts/group"
	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/egot3/supel/backend/group/internal/logctx"
	"github.com/egot3/supel/backend/group/internal/moprconv"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GroupService) AddMember(ctx context.Context, req *grpb.AddMemberRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("memberUUID", req.MemberUUID),
		slog.String("groupUUID", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	memberUUID, err := uuid.Parse(req.MemberUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	exists, err := s.Client.CheckExistance(ctx, memberUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to identify a user")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "requested user wasn't found")
	}

	subscope := "member"
	can, err := s.Client.HasPermission(ctx, ownUUID, "group", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to add a member to Group")
	}

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	is, err := s.GroupRepository.IsGroup(ctx, groupUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't determine if given uuid is assigned to a group")
	}
	if !is {
		return nil, status.Error(codes.NotFound, "couldn't find a group with this uuid")
	}

	err = s.MemberRepository.AddMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't add member to group")
	}

	return nil, nil
}

func (s *GroupService) RemoveMember(ctx context.Context, req *grpb.RemoveMemberRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("memberUUID", req.MemberUUID),
		slog.String("groupUUID", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "member"
	can, err := s.Client.HasPermission(ctx, ownUUID, "group", &subscope, rbacpb.Verb_DELETE)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to delete a member from Group")
	}

	memberUUID, err := uuid.Parse(req.MemberUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.MemberRepository.RemoveMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't delete member from group")
	}

	return nil, nil
}

func (s *GroupService) MembersGroups(ctx context.Context, req *grpb.MembersGroupsRequest) (*grpb.MembersGroupsResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("memberUUID", req.MemberUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	memberUUID, err := uuid.Parse(req.MemberUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	groupUUIDs, total, err := s.MemberRepository.MembersGroups(ctx, memberUUID, req.Page, req.Size)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, status.Error(codes.NotFound, "couldn't find member with this uuid being in groups")
		}
		return nil, status.Error(codes.Internal, "couldn't get groups for member")
	}

	groupsProto := make([]*grpb.Group, len(groupUUIDs))
	for i, groupUUID := range groupUUIDs {
		group, err := s.GroupRepository.Group(ctx, groupUUID)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't get one group for member")
		}
		groupsProto[i] = moprconv.GroupMoToPr(group)
	}

	return &grpb.MembersGroupsResponse{
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
		Groups: groupsProto,
	}, nil
}

func (s *GroupService) ListMembers(ctx context.Context, req *grpb.ListMembersRequest) (*grpb.ListMembersResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("groupUUID", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	userUUIDs, total, err := s.MemberRepository.ListMembers(ctx, groupUUID, req.Page, req.Size)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, status.Error(codes.NotFound, "couldn't find member with this uuid being in groups")
		}
		return nil, status.Error(codes.Internal, "couldn't get groups for member")
	}

	usersProto := make([]*grpb.User, len(userUUIDs))
	for i, userUUID := range userUUIDs {
		user, err := s.Client.GetUser(ctx, userUUID)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't get one member")
		}
		usersProto[i] = &grpb.User{
			Uuid:      user.Uuid,
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarUrl,
			CreatedAt: user.CreatedAt,
		}
	}

	return &grpb.ListMembersResponse{
		Members: usersProto,
		Total:   total,
		Page:    req.Page,
		Size:    req.Size,
	}, nil
}
