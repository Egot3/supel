package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	ppb "github.com/Egot3/supel/backend/contracts/puddle"
	"github.com/egot3/supel/backend/puddles/internal/carefulness"
	"github.com/egot3/supel/backend/puddles/internal/logctx"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *PuddleService) PuddleRequest(ctx context.Context, req *ppb.PuddleRequest) (*ppb.PuddleResponse, error) {
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
