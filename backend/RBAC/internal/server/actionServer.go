package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/egot3/supel/backend/rbac/internal/logctx"
	"github.com/egot3/supel/backend/rbac/types"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *RBACService) AddActions(ctx context.Context, req *rbacpb.AddActionsRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("role", req.RoleUUID),
		slog.Int("actions lenght", len(req.ActionUUIDs)),
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

	subScope := "role"
	has, err := s.RoleRepository.HasPermission(ctx, ownUUID, "rbac", &subScope, types.PATCH)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get user's abilities")
	}
	if !(uPriority > rPriority && has) {
		return nil, status.Error(codes.PermissionDenied, "don't have permissions to add actions to role")
	}

	actionUUIDs := make([]uuid.UUID, 0, len(req.ActionUUIDs))
	for _, action := range req.ActionUUIDs {
		aUUID, err := uuid.Parse(action)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "bad uuid for action")
		}
		actionUUIDs = append(actionUUIDs, aUUID)
	}

	err = s.ActionRepository.AddActions(ctx, actionUUIDs, roleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "role not found in rbac")
		}
		return nil, status.Error(codes.Internal, "couldn't add actions to role")
	}

	return nil, nil
}

func (s *RBACService) RevokeActions(ctx context.Context, req *rbacpb.RevokeActionsRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("role", req.RoleUUID),
		slog.Int("revoked actions lenght", len(req.ActionUUIDs)),
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

	subScope := "role"
	has, err := s.RoleRepository.HasPermission(ctx, ownUUID, "rbac", &subScope, types.PATCH)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get user's abilities")
	}
	if !(uPriority > rPriority && has) {
		return nil, status.Error(codes.PermissionDenied, "don't have permissions to add actions to role")
	}

	actionUUIDs := make([]uuid.UUID, 0, len(req.ActionUUIDs))
	for _, action := range req.ActionUUIDs {
		aUUID, err := uuid.Parse(action)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "bad uuid for action")
		}
		actionUUIDs = append(actionUUIDs, aUUID)
	}

	err = s.ActionRepository.RevokeActions(ctx, actionUUIDs, roleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "role not found in rbac")
		}
		return nil, status.Error(codes.Internal, "couldn't revoke actions from role")
	}

	return nil, nil
}

func (s *RBACService) ListActions(ctx context.Context, req *rbacpb.ListActionsRequest) (*rbacpb.ListActionsResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.Int("actions page", int(req.Page)),
		slog.Int("actions page's size", int(req.Size)),
	)
	ctx = logctx.WithLogger(ctx, logger)

	list, total, err := s.ActionRepository.ListActions(ctx, int(req.Page), int(req.Size))
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't list actions")
	}

	actions := make([]*rbacpb.Action, total)
	for i, action := range list {
		actions[i] = &rbacpb.Action{
			ActionUUID: action.UUID.String(),
			Scope:      action.Scope,
			Subscope:   action.SubScope,
			Verb:       rbacpb.Verb(rbacpb.Verb_value[string(action.Verb)]),
		}
	}

	return &rbacpb.ListActionsResponse{
		Action: actions,
		Total:  int64(total),
		Size:   req.Size,
		Page:   req.Page,
	}, nil
}

func (s *RBACService) ActionsByRole(ctx context.Context, req *rbacpb.ActionsByRoleRequest) (*rbacpb.ActionsByRoleResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
	)
	ctx = logctx.WithLogger(ctx, logger)

	roleUUID, err := uuid.Parse(req.GetRoleUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	actions, total, err := s.ActionRepository.ActionsByRole(ctx, roleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &rbacpb.ActionsByRoleResponse{Actions: ([]*rbacpb.Action{})}, nil
		}
		return nil, status.Error(codes.Internal, "couldn't get actions by role")
	}

	actionProto := make([]*rbacpb.Action, total)
	for i, action := range actions {
		actionProto[i] = &rbacpb.Action{
			Scope:      action.Scope,
			Subscope:   action.SubScope,
			Verb:       rbacpb.Verb(rbacpb.Verb_value[string(action.Verb)]),
			ActionUUID: action.UUID.String(),
		}
	}

	return &rbacpb.ActionsByRoleResponse{
		Actions: actionProto,
	}, nil
}

func (s *RBACService) Action(ctx context.Context, req *rbacpb.ActionRequest) (*rbacpb.ActionResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("actionUUID", req.ActionUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	actionUUID, err := uuid.Parse(req.GetActionUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	action, err := s.ActionRepository.Action(ctx, actionUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "role not found ")
		}
		return nil, status.Error(codes.Internal, "Couldn't get a role")
	}

	return &rbacpb.ActionResponse{Action: &rbacpb.Action{
		ActionUUID: actionUUID.String(),
		Scope:      action.Scope,
		Subscope:   action.SubScope,
		Verb:       rbacpb.Verb(rbacpb.Verb_value[string(action.Verb)]),
	}}, nil
}
