package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/egot3/supel/backend/rbac/internal/logctx"
	"github.com/egot3/supel/backend/rbac/types"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *RBACService) AssignorsAssignees(ctx context.Context, req *rbacpb.AssignorsAssigneesRequest) (*rbacpb.AssignorsAssigneesResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("assignor", req.AssignorUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	assignorUUID, err := uuid.Parse(req.GetAssignorUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	assignees, err := s.RoleRepository.AssignorsAssignees(ctx, assignorUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "assignor not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get assignor's assignees")
	}

	return &rbacpb.AssignorsAssigneesResponse{
		AssigneeUUIDs: assignees.Strings(),
	}, nil
}

func (s *RBACService) RoleAssignor(ctx context.Context, req *rbacpb.RoleAssignorRequest) (*rbacpb.RoleAssignorResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("assignee", req.AssigneeUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	roleUUID, err := uuid.Parse(req.GetRoleUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	assigneeUUID, err := uuid.Parse(req.GetAssigneeUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	assignorUUID, err := s.RoleRepository.RoleAssignor(ctx, roleUUID, assigneeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "assignee/role not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get assignee/role assignor")
	}

	return &rbacpb.RoleAssignorResponse{
		AssignorUUID: assignorUUID.String(),
	}, nil
}

func (s *RBACService) RoleAssignees(ctx context.Context, req *rbacpb.RoleAssigneesRequest) (*rbacpb.RoleAssigneesResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("role", req.RoleUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	roleUUID, err := uuid.Parse(req.GetRoleUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	assignees, err := s.RoleRepository.RoleAssignees(ctx, roleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "role not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get role's assignees")
	}

	return &rbacpb.RoleAssigneesResponse{
		UserUUIDs: assignees.Strings(),
	}, nil
}

func (s *RBACService) AssignRole(ctx context.Context, req *rbacpb.AssignRoleRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("role", req.RoleUUID),
		slog.String("assigneeUUID", req.AssigneeUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, _, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}
	roleUUID, err := uuid.Parse(req.GetRoleUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	rPriority, err := s.RoleRepository.RolePriority(ctx, roleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "role not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get role's priority")
	}
	uPriority, err := s.RoleRepository.UserRolePriority(ctx, ownUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get user's priority")
	}

	subScope := "role.assign"
	has, err := s.RoleRepository.HasPermission(ctx, ownUUID, "rbac", &subScope, types.POST)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get user's abilities")
	}

	if rPriority > uPriority || !has {
		return nil, status.Error(codes.PermissionDenied, "no permission to patch role")
	}

	assigneeUUID, err := uuid.Parse(req.AssigneeUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	var expiration *time.Time = nil
	if req.Expiresat != nil {
		e := req.Expiresat.AsTime()
		expiration = &e
	}

	err = s.RoleRepository.AssignRole(ctx, assigneeUUID, ownUUID, roleUUID, expiration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "role/assignee not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to assign role to assignee")
	}

	return nil, nil
}

func (s *RBACService) RevokeRole(ctx context.Context, req *rbacpb.RevokeRoleRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("role", req.RoleUUID),
		slog.String("assigneeUUID", req.AssigneeUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, _, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}
	roleUUID, err := uuid.Parse(req.GetRoleUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	rPriority, err := s.RoleRepository.RolePriority(ctx, roleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "role not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get role's priority")
	}
	uPriority, err := s.RoleRepository.UserRolePriority(ctx, ownUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get user's priority")
	}

	subScope := "role.assign"
	has, err := s.RoleRepository.HasPermission(ctx, ownUUID, "rbac", &subScope, types.DELETE)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get user's abilities")
	}

	if rPriority > uPriority || !has {
		return nil, status.Error(codes.PermissionDenied, "no permission to patch role")
	}

	assigneeUUID, err := uuid.Parse(req.AssigneeUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.RoleRepository.RevokeRole(ctx, assigneeUUID, roleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "role/assignee not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to assign role to assignee")
	}

	return nil, nil
}
