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
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *PuddleService) AddMember(ctx context.Context, req *ppb.AddMemberRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("puddleUUID", req.PuddleUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "member"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "puddle", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions to add users to puddles")
	}

	puddleUUID, err := uuid.Parse(req.PuddleUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	userUUID, err := uuid.Parse(req.UserUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad user uuid")
	}
	exists, err := s.grpcClient.CheckExistance(ctx, userUUID)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user existance")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "requested to add user does not exist")
	}

	_, err = s.puddleRepository.Puddle(ctx, puddleUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "couldn't find puddle with that uuid")
		}
		if errors.Is(err, carefulness.Gone) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Internal, "couldn't check puddle existance")
	}

	err = s.memberRepository.AddMember(ctx, userUUID, ownUUID, puddleUUID)
	if err != nil {
		if errors.Is(err, carefulness.Conflict) {
			return nil, status.Error(codes.AlreadyExists, "user to puddle connection already exists")
		}
		return nil, status.Error(codes.Internal, "couldn't add user to db")
	}

	return nil, nil
}

func (s *PuddleService) ListMembersPuddles(ctx context.Context, req *ppb.ListMembersPuddlesRequest) (*ppb.ListMembersPuddlesResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("memberUUID", req.MemberUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	memberUUID, err := uuid.Parse(req.MemberUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad user uuid")
	}

	puddleUUIDs, total, err := s.memberRepository.ListMembersPuddles(ctx, memberUUID, req.Page, req.Size)

	puddles := make([]*ppb.Puddle, len(puddleUUIDs))
	for i, puddleUUID := range puddleUUIDs {
		puddle, err := s.puddleRepository.Puddle(ctx, puddleUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, status.Error(codes.Internal, "puddle not found")
			}
			if errors.Is(err, carefulness.Gone) {
				continue
			}
			return nil, status.Error(codes.Internal, "couldn't get puddle")
		}
		memberCount, err := s.puddleRepository.PuddleMemberCount(ctx, puddleUUID)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't get puddle member count")
		}

		puddles[i] = &ppb.Puddle{
			UUID:        puddle.UUID.String(),
			Name:        puddle.Name,
			Description: puddle.Description,
			PuddleType:  ppb.PuddleType(ppb.PuddleType_value[string(puddle.PuddleType)]),

			UpdatedAt: timestamppb.New(puddle.UpdatedAt),
			CreatedAt: timestamppb.New(puddle.CreatedAt),

			MemberCount: uint64(memberCount),
		}
	}

	return &ppb.ListMembersPuddlesResponse{
		Puddles: puddles,
		Total:   uint64(total),
		Size:    req.Size,
		Page:    req.Page,
	}, nil
}

func (s *PuddleService) ListMembersPuddlesIntersections(ctx context.Context, req *ppb.ListMembersPuddlesIntersectionsRequest) (*ppb.ListMembersPuddlesIntersectionsResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("memberAUUID", req.UserAUUID),
		slog.String("memberBUUID", req.UserBUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	memberAUUID, err := uuid.Parse(req.UserAUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad user uuid")
	}

	memberBUUID, err := uuid.Parse(req.UserBUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad user uuid")
	}

	puddleUUIDs, total, err := s.memberRepository.ListMembersPuddlesIntersections(ctx, memberAUUID, memberBUUID, req.Page, req.Size)

	puddles := make([]*ppb.Puddle, len(puddleUUIDs))
	for i, puddleUUID := range puddleUUIDs {
		puddle, err := s.puddleRepository.Puddle(ctx, puddleUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, status.Error(codes.Internal, "puddle not found")
			}
			if errors.Is(err, carefulness.Gone) {
				continue
			}
			return nil, status.Error(codes.Internal, "couldn't get puddle")
		}
		memberCount, err := s.puddleRepository.PuddleMemberCount(ctx, puddleUUID)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't get puddle member count")
		}

		puddles[i] = &ppb.Puddle{
			UUID:        puddle.UUID.String(),
			Name:        puddle.Name,
			Description: puddle.Description,
			PuddleType:  ppb.PuddleType(ppb.PuddleType_value[string(puddle.PuddleType)]),

			UpdatedAt: timestamppb.New(puddle.UpdatedAt),
			CreatedAt: timestamppb.New(puddle.CreatedAt),

			MemberCount: uint64(memberCount),
		}
	}

	return &ppb.ListMembersPuddlesIntersectionsResponse{
		Puddles: puddles,
		Total:   uint64(total),
		Size:    req.Size,
		Page:    req.Page,
	}, nil
}

func (s *PuddleService) ListAddersAddors(ctx context.Context, req *ppb.ListAddersAddorsRequest) (*ppb.ListAddersAddorsResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("adderUUID", req.AdderUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	adderUUID, err := uuid.Parse(req.AdderUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad user uuid")
	}

	memberUUIDs, total, err := s.memberRepository.ListAddersAddors(ctx, adderUUID, req.Page, req.Size)

	members := make([]*ppb.User, len(memberUUIDs))
	for i, memberUUID := range memberUUIDs {
		member, err := s.grpcClient.GetUser(ctx, memberUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, status.Error(codes.Internal, "puddle not found")
			}
			if errors.Is(err, carefulness.Gone) {
				continue
			}
			return nil, status.Error(codes.Internal, "couldn't get puddle")
		}
		isMod, err := s.moderatorRepository.IsModerator(ctx, memberUUID, memberUUID)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't get puddle member count")
		}

		members[i] = &ppb.User{
			Uuid:        memberUUID.String(),
			Nickname:    member.Nickname,
			AvatarUrl:   member.AvatarUrl,
			IsModerator: isMod,
			CreatedAt:   member.CreatedAt,
		}
	}

	return &ppb.ListAddersAddorsResponse{
		Users: members,
		Total: uint64(total),
		Size:  req.Size,
		Page:  req.Page,
	}, nil
}
