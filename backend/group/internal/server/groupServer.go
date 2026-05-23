package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	grpb "github.com/Egot3/supel/backend/contracts/group"
	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/egot3/supel/backend/group/internal/carefulness"
	"github.com/egot3/supel/backend/group/internal/logctx"
	"github.com/egot3/supel/backend/group/internal/models"
	"github.com/egot3/supel/backend/group/internal/moprconv"
	"github.com/egot3/supel/backend/group/internal/types"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GroupService) Group(ctx context.Context, req *grpb.GroupRequest) (*grpb.GroupResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("group_uuid", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	groupUUID, err := uuid.Parse(req.GetGroupUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	group, err := s.GroupRepository.Group(ctx, groupUUID)
	if err != nil {
		if errors.Is(err, carefulness.Gone) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.NotFound, "group with this uuid is gone")
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "group with this uuid was not found")
		}
		return nil, status.Error(codes.Internal, "couldn't get group with that uuid")
	}

	return &grpb.GroupResponse{
		Group: moprconv.GroupMoToPr(group),
	}, nil
}

func (s *GroupService) CreateGroup(ctx context.Context, req *grpb.CreateGroupRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("name", req.Name),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	can, err := s.Client.HasPermission(ctx, ownUUID, "group", nil, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to createGroup")
	}

	curatorUUID, err := uuid.Parse(req.CuratorUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	if req.GroupType.String() == types.GROUP {
		is, err := s.CuratorRepository.IsCurator(ctx, curatorUUID)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't check if user is curator")
		}
		if !is {
			return nil, status.Error(codes.PermissionDenied, "can't create a learning group if not a curator")
		}
	}

	err = s.GroupRepository.CreateGroup(ctx, curatorUUID, req.Name, req.Description, types.GroupType(req.GroupType.String()))
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't create group")
	}

	return nil, nil
}

func (s *GroupService) SearchGroup(ctx context.Context, req *grpb.SearchGroupRequest) (*grpb.SearchGroupResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("sample", req.Sample),
		slog.Int("limit", int(req.GetLimit())),
	)
	ctx = logctx.WithLogger(ctx, logger)

	serp, err := s.GroupRepository.Search(ctx, req.Sample, int(req.Limit))
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, status.Error(codes.NotFound, "couldn't find anything for this query")
		}
		return nil, status.Error(codes.Internal, "couldn't make a search")
	}

	serpProto := make([]*grpb.Group, len(serp))
	for i, result := range serp {
		serpProto[i] = moprconv.GroupMoToPr(&result)
	}

	return &grpb.SearchGroupResponse{
		Groups: serpProto,
	}, nil
}

func (s *GroupService) DeleteGroup(ctx context.Context, req *grpb.DeleteGroupRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("group_uuid", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	can, err := s.Client.HasPermission(ctx, ownUUID,
		"group",
		nil,
		rbacpb.Verb_DELETE)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to delete a Group")
	}

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.GroupRepository.DeleteGroup(ctx, groupUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't delete group")
	}

	return nil, nil
}

func (s *GroupService) ListGroups(ctx context.Context, req *grpb.ListGroupsRequest) (*grpb.ListGroupsResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.Int("page", int(req.Page)),
		slog.Int("size", int(req.Size)),
	)
	ctx = logctx.WithLogger(ctx, logger)

	list, total, err := s.GroupRepository.ListGroups(ctx, req.Page, req.Size, types.GroupType(req.GroupType), bun.OrderDesc)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to get groups for listing")
	}

	listProto := make([]*grpb.Group, len(list))
	for i, group := range list {
		listProto[i] = moprconv.GroupMoToPr(&group)
	}

	return &grpb.ListGroupsResponse{
		Groups: listProto,
		Total:  total,
		Page:   req.Page,
		Size:   req.Size,
	}, nil
}

func (s *GroupService) CuratorsGroups(ctx context.Context, req *grpb.CuratorsGroupsRequest) (*grpb.CuratorsGroupsResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("curatorUUID", req.CuratorUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	curatorUUID, err := uuid.Parse(req.CuratorUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	groups, err := s.GroupRepository.CuratorsGroups(ctx, curatorUUID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, status.Error(codes.NotFound, "couldn't find curator with this uuid owning groups")
		}
		return nil, status.Error(codes.Internal, "couldn't get groups for curator")
	}

	groupsProto := make([]*grpb.Group, len(groups))
	for i, group := range groups {
		groupsProto[i] = moprconv.GroupMoToPr(&group)
	}

	return &grpb.CuratorsGroupsResponse{
		Groups: groupsProto,
	}, nil
}

func (s *GroupService) PatchGroup(ctx context.Context, req *grpb.PatchGroupRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("groupUUID", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	can, err := s.Client.HasPermission(ctx, ownUUID, "group", nil, rbacpb.Verb_PATCH)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to patch group")
	}

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	canEdit, err := s.CuratorRepository.CanEdit(ctx, ownUUID, groupUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to check curator's group relation")
	}
	if !canEdit {
		return nil, status.Error(codes.PermissionDenied, "can't edit this group")
	}

	err = s.GroupRepository.PatchGroup(ctx,
		models.GroupPatched{
			UUID:        groupUUID,
			Name:        req.Name,
			Description: req.Description,
			GroupType:   types.GroupType(req.GroupType.String())})

	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't patch group")
	}

	return nil, nil
}
