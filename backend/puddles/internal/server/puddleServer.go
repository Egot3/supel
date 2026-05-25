package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	ppb "github.com/Egot3/supel/backend/contracts/puddle"
	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/egot3/supel/backend/puddles/internal/carefulness"
	"github.com/egot3/supel/backend/puddles/internal/logctx"
	"github.com/egot3/supel/backend/puddles/internal/models"
	"github.com/egot3/supel/backend/puddles/internal/types"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *PuddleService) Puddle(ctx context.Context, req *ppb.PuddleRequest) (*ppb.PuddleResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("puddleUUID", req.PuddleUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	puddleUUID, err := uuid.Parse(req.PuddleUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	puddle, err := s.puddleRepository.Puddle(ctx, puddleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested puddle doesn't seem to exist")
		}
		if errors.Is(err, carefulness.Gone) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Internal, "couldn't get puddle")
	}
	memberCount, err := s.puddleRepository.PuddleMemberCount(ctx, puddleUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't get puddle member count")
	}

	return &ppb.PuddleResponse{
		Puddle: &ppb.Puddle{
			UUID:        puddle.UUID.String(),
			Name:        puddle.Name,
			Description: puddle.Description,
			PuddleType:  ppb.PuddleType(ppb.PuddleType_value[string(puddle.PuddleType)]),

			UpdatedAt: timestamppb.New(puddle.UpdatedAt),
			CreatedAt: timestamppb.New(puddle.CreatedAt),

			MemberCount: uint64(memberCount),
		},
	}, nil
}

func (s *PuddleService) CreatePuddle(ctx context.Context, req *ppb.CreatePuddleRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("puddleName", req.Name),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "puddle", nil, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for puddle creation")
	}

	if len(req.Name) <= 3 {
		return nil, status.Error(codes.InvalidArgument, "puddle's name must be > than 3 chars")
	}

	startingUserUUIDs := make(uuid.UUIDs, len(req.InvitedUsersUUIDs))
	for i, user := range req.InvitedUsersUUIDs {
		startingUserUUIDs[i], err = uuid.Parse(user)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Bad invited user uuid")
		}

		exists, err := s.grpcClient.CheckExistance(ctx, startingUserUUIDs[i])
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't check if invited user exists")
		}
		if !exists {
			return nil, status.Error(codes.NotFound, "invited user wasn't found")
		}
	}

	switch types.PuddleTypePrToTy(req.PuddleType) {
	case types.ONEONONE:
		if len(startingUserUUIDs) > 2 {
			return nil, status.Error(codes.InvalidArgument, "can't create one on one puddle with more than two users")
		}
		if len(startingUserUUIDs) < 2 {
			return nil, status.Error(codes.InvalidArgument, "can't create on on one puddle with less than two users")
		}
		err = s.puddleRepository.CreateOneOnOnePuddle(ctx, req.Name, req.Description, [2]uuid.UUID(startingUserUUIDs))
	case types.GROUP:
		err = s.puddleRepository.CreateGroup(ctx, req.Name, req.Description, startingUserUUIDs, ownUUID)
	case types.CHANNEL:
		err = s.puddleRepository.CreateChannel(ctx, req.Name, req.Description, ownUUID)
	default:
		return nil, status.Error(codes.InvalidArgument, "can't create group with no type")
	}

	if err != nil {
		if errors.Is(err, carefulness.Conflict) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.AlreadyExists, "such puddle already exists")
		}
		return nil, status.Error(codes.Internal, "couldn't create puddle")
	}

	return nil, nil
}

func (s *PuddleService) PatchPuddle(ctx context.Context, req *ppb.PatchPuddleRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("puddleUUID", req.PuddleUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "puddle", nil, rbacpb.Verb_PATCH)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for puddle patching")
	}

	puddleUUID, err := uuid.Parse(req.PuddleUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	is, err := s.moderatorRepository.IsModerator(ctx, puddleUUID, ownUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check if user is a moderator")
	}
	if !is {
		return nil, status.Error(codes.PermissionDenied, "can't patch puddle when not a moderator")
	}

	if req.Name != nil && len(*req.Name) <= 3 {
		return nil, status.Error(codes.InvalidArgument, "puddle's name must be > than 3 chars")
	}

	err = s.puddleRepository.PatchPuddle(ctx, models.PuddlePatched{
		UUID:        puddleUUID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "puddle with this uuid wasn't found")
		}
		return nil, status.Error(codes.Internal, "couldn't patch puddle")
	}

	return nil, nil
}

func (s *PuddleService) DeletePuddle(ctx context.Context, req *ppb.DeletePuddleRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("puddleUUID", req.PuddleUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "puddle", nil, rbacpb.Verb_DELETE)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for puddle deletion")
	}

	puddleUUID, err := uuid.Parse(req.PuddleUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	is, err := s.moderatorRepository.IsModerator(ctx, puddleUUID, ownUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check if user is a moderator")
	}
	if !is {
		return nil, status.Error(codes.PermissionDenied, "can't delete puddle when not a moderator")
	}

	err = s.puddleRepository.DeletePuddle(ctx, puddleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "couldn't find requested puddle")
		}
		return nil, status.Error(codes.Internal, "couldn't delete puddle")
	}

	return nil, nil
}

func (s *PuddleService) ListPuddleMembers(ctx context.Context, req *ppb.ListPuddleMembersRequest) (*ppb.ListPuddleMembersResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("puddleUUID", req.PuddleUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	if req.Size == 0 {
		return nil, status.Error(codes.InvalidArgument, "size can't be zero")
	}

	puddleUUID, err := uuid.Parse(req.PuddleUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	members, total, err := s.puddleRepository.ListPuddleMembers(ctx, puddleUUID, req.Page, req.Size)

	membersProto := make([]*ppb.User, len(members))
	for i, member := range members {
		user, err := s.grpcClient.GetUser(ctx, member)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't fetch member's info")
		}
		isMod, err := s.moderatorRepository.IsModerator(ctx, puddleUUID, member)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't check if user is a moderator")
		}

		membersProto[i] = &ppb.User{
			Uuid:        user.Uuid,
			Nickname:    user.Nickname,
			AvatarUrl:   user.AvatarUrl,
			CreatedAt:   user.CreatedAt,
			IsModerator: isMod,
		}
	}

	return &ppb.ListPuddleMembersResponse{
		Users: membersProto,
		Total: uint64(total),
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}
