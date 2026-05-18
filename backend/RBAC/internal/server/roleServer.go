package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/egot3/supel/backend/rbac/internal/logctx"
	"github.com/egot3/supel/backend/rbac/internal/models"
	"github.com/egot3/supel/backend/rbac/types"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *RBACService) HasPermission(ctx context.Context, req *rbacpb.HasPermissionQuestion) (*rbacpb.HasPermissionAnswer, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("user_uuid", req.UserUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	userUUID, err := uuid.Parse(req.GetUserUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	has, err := s.RoleRepository.HasPermission(ctx, userUUID, req.GetScope(), req.SubScope, types.Verb(req.Verb.String()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Not found for permissions")
		}
		return nil, status.Error(codes.Internal, "Couldn't check if user has perms")
	}

	return &rbacpb.HasPermissionAnswer{Has: has}, nil
}

func (s *RBACService) DeleteRole(ctx context.Context, req *rbacpb.DeleteRoleRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("role to delete", req.RoleUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, _, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subScope := "role"
	has, err := s.RoleRepository.HasPermission(ctx, ownUUID, "rbac", &subScope, types.DELETE)
	if err != nil || !has {
		return nil, status.Error(codes.PermissionDenied, "no permission to delete role")
	}

	roleUUID, err := uuid.Parse(req.GetRoleUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.RoleRepository.DeleteRole(ctx, roleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Not found role for deletion")
		}
		return nil, status.Error(codes.Internal, "Couldn't delete role")
	}

	return nil, nil
}

func (s *RBACService) Role(ctx context.Context, req *rbacpb.RoleRequest) (*rbacpb.RoleResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("role to get", req.RoleUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	roleUUID, err := uuid.Parse(req.GetRoleUUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	role, err := s.RoleRepository.Role(ctx, roleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "role not found ")
		}
		return nil, status.Error(codes.Internal, "Couldn't get a role")
	}

	var extended *string = nil
	if role.ExtendedUUID != nil {
		e := role.ExtendedUUID.String()
		extended = &e
	}
	var deletedAt *timestamppb.Timestamp = nil
	if d := role.DeletedAt; d != nil {
		deletedAt = timestamppb.New(*d)
	}
	return &rbacpb.RoleResponse{
		RoleUUID:     role.UUID.String(),
		Name:         role.Name,
		Description:  role.Description,
		ExtendedUUID: extended,
		Priority:     int32(role.Priority),
		Createdat:    timestamppb.New(role.CreatedAt),
		Deletedat:    deletedAt,
		Updatedat:    timestamppb.New(role.UpdatedAt),
	}, nil
}

func (s *RBACService) PatchRole(ctx context.Context, req *rbacpb.PatchRoleRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("role to patch", req.RoleUUID),
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

	if !(rPriority < uPriority && has) {
		return nil, status.Error(codes.PermissionDenied, "no permission to patch role")
	}
	var priority *int16 = nil
	if p := req.Priority; p != nil {
		pv := *p
		pv16 := int16(pv)
		priority = &pv16
	}
	err = s.RoleRepository.PatchRole(ctx, models.PatchedRole{
		UUID:         roleUUID,
		Name:         req.Name,
		Description:  req.Description,
		ExtendedUUID: req.ExtendedUUID,
		Priority:     priority,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "role not found in rbac")
		}
		return nil, status.Error(codes.Internal, "couldn't patch a role")
	}

	return nil, nil
}

func (s *RBACService) CreateRole(ctx context.Context, req *rbacpb.CreateRoleRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("role to create", req.Name),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, _, ok := UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	uPriority, err := s.RoleRepository.UserRolePriority(ctx, ownUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get user's priority")
	}

	subScope := "role"
	has, err := s.RoleRepository.HasPermission(ctx, ownUUID, "rbac", &subScope, types.POST)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found in rbac")
		}
		return nil, status.Error(codes.Internal, "failed to get user's abilities")
	}

	if !(int16(req.Priority) < uPriority && has) {
		return nil, status.Error(codes.PermissionDenied, "no permission to create role")
	}

	var extendedUUID *uuid.UUID = nil
	if eus := req.ExtendedUUID; eus != nil {
		eu, err := uuid.Parse(*eus)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Bad uuid")
		}
		extendedUUID = &eu
	}

	actionUUIDs := make([]uuid.UUID, 0, len(req.Actions))
	for _, action := range req.Actions {
		aUUID, err := uuid.Parse(action.ActionUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("bad uuid for action %v.%v:%v", action.Scope, action.Subscope, action.Verb.String()))
		}
		actionUUIDs = append(actionUUIDs, aUUID)
	}

	err = s.RoleRepository.CreateRole(ctx, models.Role{
		Name:         req.Name,
		Description:  req.Description,
		ExtendedUUID: extendedUUID,
		Priority:     int16(req.Priority),
	}, actionUUIDs)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't create role")
	}

	return nil, nil
}

func (s *RBACService) ListRoles(ctx context.Context, req *rbacpb.ListRolesRequest) (*rbacpb.ListRolesResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.Int("roles page", int(req.Page)),
		slog.Int("roles page size", int(req.Size)),
	)
	ctx = logctx.WithLogger(ctx, logger)

	list, total, err := s.RoleRepository.ListRoles(ctx, int(req.Page), int(req.Size))
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to get roles for listing")
	}

	roles := make([]*rbacpb.ListRolesResponse_Role, total)
	for i, role := range list {
		var extendedUUID *string = nil
		if eu := role.ExtendedUUID; eu != nil {
			eus := eu.String()
			extendedUUID = &(eus)
		}
		var deletedAt *timestamppb.Timestamp = nil
		if d := role.DeletedAt; d != nil {
			deletedAt = timestamppb.New(*d)
		}
		roles[i] = &rbacpb.ListRolesResponse_Role{
			RoleUUID:     role.UUID.String(),
			Name:         role.Name,
			Description:  role.Description,
			ExtendedUUID: extendedUUID,
			Priority:     int32(role.Priority),
			Createdat:    timestamppb.New(role.CreatedAt),
			Deletedat:    deletedAt,
			Updatedat:    timestamppb.New(role.UpdatedAt),
		}
	}

	return &rbacpb.ListRolesResponse{
		Roles: roles,
		Total: int64(total),
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}
