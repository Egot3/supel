package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"github.com/egot3/supel/backend/timetable/internal/carefulness"
	"github.com/egot3/supel/backend/timetable/internal/logctx"
	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *TimetableService) Period(ctx context.Context, req *ttpb.PeriodRequest) (*ttpb.PeriodResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("periodUUID", req.PeriodUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	periodUUID, err := uuid.Parse(req.PeriodUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	period, err := s.periodRepository.Period(ctx, periodUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested period wasn't found")
		}
		if errors.Is(err, carefulness.Gone) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Internal, "couldn't get period")
	}

	return &ttpb.PeriodResponse{
		Period: &ttpb.Period{
			UUID:      period.UUID.String(),
			Name:      period.Name,
			Position:  int32(period.Position),
			StartTime: timestamppb.New(period.StartTime),
			EndTime:   timestamppb.New(period.EndTime),

			UpdatedAt: timestamppb.New(period.UpdatedAt),
			CreatedAt: timestamppb.New(period.CreatedAt),
		},
	}, nil
}

func (s *TimetableService) CreatePeriod(ctx context.Context, req *ttpb.CreatePeriodRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("name", req.Name),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "period"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for puddle creation")
	}

	err = s.periodRepository.CreatePeriod(ctx, req.Name, int16(req.Position), req.StartTime.AsTime(), req.StartTime.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't create request")
	}

	return nil, nil
}

func (s *TimetableService) PatchPeriod(ctx context.Context, req *ttpb.PatchPeriodRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("periodUUID", req.PeriodUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "period"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for puddle creation")
	}

	periodUUID, err := uuid.Parse(req.PeriodUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	var startTimePtr *time.Time = nil
	if req.StartTime != nil {
		startTime := req.StartTime.AsTime()
		startTimePtr = &startTime
	}

	var endTimePtr *time.Time = nil
	if req.StartTime != nil {
		endTime := req.StartTime.AsTime()
		startTimePtr = &endTime
	}

	err = s.periodRepository.PatchPeriod(ctx, models.PeriodPatched{
		UUID:      periodUUID,
		Name:      req.Name,
		Position:  req.Position,
		StartTime: startTimePtr,
		EndTime:   endTimePtr,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "period to patch wasn't found")
		}
		return nil, status.Error(codes.Internal, "couldn't patch period")
	}

	return nil, nil
}

func (s *TimetableService) DeletePeriod(ctx context.Context, req *ttpb.DeletePeriodRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("periodUUID", req.PeriodUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "period"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_DELETE)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for puddle creation")
	}

	periodUUID, err := uuid.Parse(req.PeriodUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.periodRepository.DeletePeriod(ctx, periodUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "period to delete wasn't found")
		}
		return nil, status.Error(codes.Internal, "couldn't delete period")
	}

	return nil, nil
}

func (s *TimetableService) ListPeriods(ctx context.Context, req *ttpb.ListPeriodsRequest) (*ttpb.ListPeriodsResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.Int("size", int(req.Size)),
		slog.Int("page", int(req.Page)),
	)
	ctx = logctx.WithLogger(ctx, logger)

	periods, total, err := s.periodRepository.ListPeriods(ctx, req.Page, req.Size)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't list periods")
	}

	periodsProto := make([]*ttpb.Period, len(periods))
	for i, period := range periods {
		periodsProto[i] = &ttpb.Period{
			UUID:      period.UUID.String(),
			Name:      period.Name,
			Position:  int32(period.Position),
			StartTime: timestamppb.New(period.StartTime),
			EndTime:   timestamppb.New(period.EndTime),

			UpdatedAt: timestamppb.New(period.UpdatedAt),
			CreatedAt: timestamppb.New(period.CreatedAt),
		}
	}

	return &ttpb.ListPeriodsResponse{
		Page:    req.Page,
		Size:    req.Size,
		Total:   uint64(total),
		Periods: periodsProto,
	}, nil
}
