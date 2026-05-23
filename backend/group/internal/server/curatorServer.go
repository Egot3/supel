package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	grpb "github.com/Egot3/supel/backend/contracts/group"
	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/egot3/supel/backend/group/internal/logctx"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *GroupService) AssignCuratorToSenior(ctx context.Context, req *grpb.AssignCuratorToSeniorRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("subordinateUUID", req.SubordinateUUID),
		slog.String("seniorUUID", req.SeniorUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "curator.heirarchy"
	can, err := s.Client.HasPermission(ctx,
		ownUUID,
		"group",
		&subscope,
		rbacpb.Verb_POST,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to assign a sub to sen")
	}

	subordinateUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get subordinate uuid from token")
	}
	exists, err := s.Client.CheckExistance(ctx, subordinateUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to identify a user")
	}
	if exists {
		return nil, status.Error(codes.NotFound, "requested user wasn't found")
	}

	is, err := s.CuratorRepository.IsCurator(ctx, subordinateUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to find sub in curator repository")
	}
	if !is {
		return nil, status.Error(codes.PermissionDenied, "can't assign not a curator to sen")
	}

	seniorUUID, err := uuid.Parse(req.SeniorUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	is, err = s.CuratorRepository.IsCurator(ctx, seniorUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to find sen in curator repository")
	}
	if !is {
		return nil, status.Error(codes.PermissionDenied, "can't assign not a curator to sub")
	}

	will, err := s.CuratorRepository.WillCycle(ctx, seniorUUID, subordinateUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check for cycles")
	}
	if will {
		return nil, status.Error(codes.FailedPrecondition, "this assignment will create a cycle, terminating")
	}

	err = s.CuratorRepository.AssignCuratorToSenior(ctx, seniorUUID, subordinateUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't assign sub to sen")
	}

	return nil, nil
}

func (s *GroupService) AssignCuratorToGroup(ctx context.Context, req *grpb.AssignCuratorToGroupRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("curatorUUID", req.CuratorUUID),
		slog.String("groupUUID", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "group.curator"
	can, err := s.Client.HasPermission(ctx, ownUUID, "group", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to assign a cur to gr")
	}

	curUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get cur uuid from token")
	}
	exist, err := s.Client.CheckExistance(ctx, curUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to identify a user")
	}
	if exist {
		return nil, status.Error(codes.NotFound, "requested user wasn't found")
	}

	is, err := s.CuratorRepository.IsCurator(ctx, curUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to find sub in curator repository")
	}
	if !is {
		return nil, status.Error(codes.PermissionDenied, "can't assign not a curator to sen")
	}

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.CuratorRepository.AssignCuratorToGroup(ctx, curUUID, groupUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't assign cur to gr")
	}

	return nil, nil
}

func (s *GroupService) RevokeCuratorFromSenior(ctx context.Context, req *grpb.RevokeCuratorFromSeniorRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("subordinateUUID", req.SubordinateUUID),
		slog.String("seniorUUID", req.SeniorUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "curator.heirarchy"
	can, err := s.Client.HasPermission(ctx, ownUUID,
		"group",
		&subscope,
		rbacpb.Verb_DELETE,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to revoke a sub from sen")
	}

	subordinateUUID, err := uuid.Parse(req.SubordinateUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	seniorUUID, err := uuid.Parse(req.SeniorUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.CuratorRepository.RevokeCuratorFromSenior(ctx, seniorUUID, subordinateUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't revoke sub from sen")
	}

	return nil, nil
}

func (s *GroupService) RevokeCuratorFromGroup(ctx context.Context, req *grpb.RevokeCuratorFromGroupRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("curatorUUID", req.CuratorUUID),
		slog.String("groupUUID", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "group.curator"
	can, err := s.Client.HasPermission(ctx, ownUUID, "group", &subscope, rbacpb.Verb_DELETE)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to revoke a cur from gr")
	}

	curatorUUID, err := uuid.Parse(req.CuratorUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.CuratorRepository.RevokeCuratorFromGroup(ctx, curatorUUID, groupUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't revoke cur from gr")
	}

	return nil, nil
}

func (s *GroupService) AddCurator(ctx context.Context, req *grpb.AddCuratorRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("curatorUUID", req.CuratorUUID),
		slog.String("seniorUUID", req.SeniorUUID),
		slog.String("groupUUID", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "curator"
	can, err := s.Client.HasPermission(ctx, ownUUID, "group", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to create a cur")
	}

	curatorUUID, err := uuid.Parse(req.CuratorUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	exists, err := s.Client.CheckExistance(ctx, curatorUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to identify a user")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "requested user wasn't found")
	}

	is, err := s.CuratorRepository.IsCurator(ctx, curatorUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to find sub in curator repository")
	}
	if is {
		return nil, status.Error(codes.AlreadyExists, "can't recreate curator")
	}

	seniorUUID, err := uuid.Parse(req.SeniorUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	is, err = s.CuratorRepository.IsCurator(ctx, seniorUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to find sen in curator repository")
	}
	if !is {
		return nil, status.Error(codes.PermissionDenied, "can't assign not a curator to sub")
	}

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.CuratorRepository.AddCurator(ctx, seniorUUID, curatorUUID, groupUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't create a cur")
	}

	return nil, nil
}

func (s *GroupService) RevokeCurator(ctx context.Context, req *grpb.RevokeCuratorRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("curatorUUID", req.CuratorUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "curator"
	can, err := s.Client.HasPermission(ctx, ownUUID, "group", &subscope, rbacpb.Verb_DELETE)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to delete a cur")
	}

	curatorUUID, err := uuid.Parse(req.CuratorUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.CuratorRepository.RevokeCurator(ctx, curatorUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user is not a cur")
		}
		return nil, status.Error(codes.Internal, "couldn't delete a cur")
	}

	return nil, nil
}

func (s *GroupService) GroupsCurators(ctx context.Context, req *grpb.GroupsCuratorsRequest) (*grpb.GroupsCuratorsResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("groupUUID", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	curatorUUIDs, err := s.CuratorRepository.GroupsCurators(ctx, groupUUID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, status.Error(codes.NotFound, "couldn't find member with this uuid being in groups")
		}
		return nil, status.Error(codes.Internal, "couldn't get groups for member")
	}

	curatorsProto := make([]*grpb.User, len(curatorUUIDs))
	for i, curatorUUID := range curatorUUIDs {
		user, err := s.Client.GetUser(ctx, curatorUUID)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't get one member")
		}
		curatorsProto[i] = &grpb.User{
			Uuid:      user.Uuid,
			Nickname:  user.Nickname,
			AvatarUrl: user.AvatarUrl,
			CreatedAt: user.CreatedAt,
		}
	}

	return &grpb.GroupsCuratorsResponse{
		Curators: curatorsProto,
	}, nil
}
