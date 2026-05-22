package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	pb "github.com/Egot3/supel/backend/contracts"
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
	can, err := s.RBACClient.HasPermission(ctx, &rbacpb.HasPermissionQuestion{
		Scope:    "group",
		SubScope: &subscope,
		Verb:     rbacpb.Verb_POST,
		UserUUID: ownUUID.String(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can.Has {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to assign a sub to sen")
	}

	identityResp, err := s.IdentityClient.CheckExistance(ctx, &pb.CheckExistanceRequest{Uuid: req.SubordinateUUID})
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to identify a user")
	}
	if !identityResp.Exists {
		return nil, status.Error(codes.NotFound, "requested user wasn't found")
	}

	subordinateUUID, err := uuid.Parse(req.SubordinateUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
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
	can, err := s.RBACClient.HasPermission(ctx, &rbacpb.HasPermissionQuestion{
		Scope:    "group",
		SubScope: &subscope,
		Verb:     rbacpb.Verb_POST,
		UserUUID: ownUUID.String(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can.Has {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to assign a cur to gr")
	}

	identityResp, err := s.IdentityClient.CheckExistance(ctx, &pb.CheckExistanceRequest{Uuid: req.CuratorUUID})
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to identify a user")
	}
	if !identityResp.Exists {
		return nil, status.Error(codes.NotFound, "requested user wasn't found")
	}

	curatorUUID, err := uuid.Parse(req.CuratorUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	is, err := s.CuratorRepository.IsCurator(ctx, curatorUUID)
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

	err = s.CuratorRepository.AssignCuratorToGroup(ctx, curatorUUID, groupUUID)
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
	can, err := s.RBACClient.HasPermission(ctx, &rbacpb.HasPermissionQuestion{
		Scope:    "group",
		SubScope: &subscope,
		Verb:     rbacpb.Verb_DELETE,
		UserUUID: ownUUID.String(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can.Has {
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
	can, err := s.RBACClient.HasPermission(ctx, &rbacpb.HasPermissionQuestion{
		Scope:    "group",
		SubScope: &subscope,
		Verb:     rbacpb.Verb_DELETE,
		UserUUID: ownUUID.String(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can.Has {
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
	can, err := s.RBACClient.HasPermission(ctx, &rbacpb.HasPermissionQuestion{
		Scope:    "group",
		SubScope: &subscope,
		Verb:     rbacpb.Verb_POST,
		UserUUID: ownUUID.String(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can.Has {
		return nil, status.Error(codes.PermissionDenied, "don't have enough permissions to create a cur")
	}

	identityResp, err := s.IdentityClient.CheckExistance(ctx, &pb.CheckExistanceRequest{Uuid: req.CuratorUUID})
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to identify a user")
	}
	if !identityResp.Exists {
		return nil, status.Error(codes.NotFound, "requested user wasn't found")
	}

	curatorUUID, err := uuid.Parse(req.CuratorUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
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
	can, err := s.RBACClient.HasPermission(ctx, &rbacpb.HasPermissionQuestion{
		Scope:    "group",
		SubScope: &subscope,
		Verb:     rbacpb.Verb_DELETE,
		UserUUID: ownUUID.String(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't call an rbac to check permissions")
	}
	if !can.Has {
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
		user, err := s.UserClient.GetUser(ctx, &pb.GetUserRequest{Uuid: curatorUUID.String()})
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't get one member")
		}
		curatorsProto[i] = &grpb.User{
			Uuid:      user.User.Uuid,
			Nickname:  user.User.Nickname,
			AvatarUrl: user.User.AvatarUrl,
			CreatedAt: user.User.CreatedAt,
		}
	}

	return &grpb.GroupsCuratorsResponse{
		Curators: curatorsProto,
	}, nil
}
