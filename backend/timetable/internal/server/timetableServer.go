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

func (s *TimetableService) Timetable(ctx context.Context, req *ttpb.TimetableRequest) (*ttpb.TimetableResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("timetableUUID", req.TimetableUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	timetableUUID, err := uuid.Parse(req.TimetableUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	timetable, err := s.timetableRepository.Timetable(ctx, timetableUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested timetable wasn't found")
		}
		if errors.Is(err, carefulness.Gone) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Internal, "couldn't get timetable")
	}

	var revokeAt *timestamppb.Timestamp = nil
	if timetable.RevokeAt != nil {
		ra := *timetable.RevokeAt
		revokeAt = timestamppb.New(ra)
	}

	var assignAt *timestamppb.Timestamp = nil
	if timetable.AssignAt != nil {
		aa := *timetable.AssignAt
		revokeAt = timestamppb.New(aa)
	}
	return &ttpb.TimetableResponse{
		Timetable: &ttpb.Timetable{
			UUID:      timetable.UUID.String(),
			GroupUUID: timetable.GroupUUID.String(),
			Name:      timetable.Name,
			AssignAt:  assignAt,
			RevokeAt:  revokeAt,
		},
	}, nil
}

func (s *TimetableService) TimetableByDate(ctx context.Context, req *ttpb.TimetableByDateRequest) (*ttpb.TimetableByDateResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("date", req.Date),
		slog.String("groupUUID", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	dateTime, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "bad date")
	}

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	timetable, err := s.timetableRepository.TimetableByDate(ctx, groupUUID, dateTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Timetable at this time was't found")
		}
		return nil, status.Error(codes.Internal, "couldn't get timetable by date")
	}

	var revokeAt *timestamppb.Timestamp = nil
	if timetable.RevokeAt != nil {
		ra := *timetable.RevokeAt
		revokeAt = timestamppb.New(ra)
	}

	var assignAt *timestamppb.Timestamp = nil
	if timetable.AssignAt != nil {
		aa := *timetable.AssignAt
		revokeAt = timestamppb.New(aa)
	}
	return &ttpb.TimetableByDateResponse{
		Timetable: &ttpb.Timetable{
			UUID:      timetable.UUID.String(),
			GroupUUID: timetable.GroupUUID.String(),
			Name:      timetable.Name,
			AssignAt:  assignAt,
			RevokeAt:  revokeAt,
		},
	}, nil
}

func (s *TimetableService) CreateTimetable(ctx context.Context, req *ttpb.CreateTimetableRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("groupUUID", req.GroupUUID),
		slog.String("name", req.Name),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", nil, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for timetable creation")
	}

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	var revokeAt *time.Time = nil
	if req.RevokeAt != nil {
		ra := req.RevokeAt.AsTime()
		revokeAt = &ra
	}

	var assignAt *time.Time = nil
	if req.AssignAt != nil {
		aa := req.AssignAt.AsTime()
		revokeAt = &aa
	}
	err = s.timetableRepository.CreateTimetable(ctx, groupUUID, req.Name, assignAt, revokeAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "coudln't create timetable")
	}

	return nil, nil
}

func (s *TimetableService) PatchTimetable(ctx context.Context, req *ttpb.PatchTimetableRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("timetableUUID", req.TimetableUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", nil, rbacpb.Verb_PATCH)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for timetable patching")
	}

	var groupUUID uuid.UUID = uuid.Nil
	if req.GroupUUID != nil {
		groupUUID, err = uuid.Parse(*req.GroupUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Bad uuid")
		}
	}
	timetableUUID, err := uuid.Parse(req.TimetableUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	var revokeAt *time.Time = nil
	revokeAtUpdated := false
	if req.RevokeAt != nil {
		ra := req.RevokeAt.AsTime()
		revokeAt = &ra
		revokeAtUpdated = true
	}

	var assignAt *time.Time = nil
	assignAtUpdated := false
	if req.AssignAt != nil {
		aa := req.AssignAt.AsTime()
		revokeAt = &aa
		assignAtUpdated = true
	}

	err = s.timetableRepository.PatchTimetable(ctx, models.TimetablePatched{
		UUID:             timetableUUID,
		GroupUUID:        groupUUID,
		Name:             req.Name,
		AssignAtUpdated:  assignAtUpdated,
		AssignAt:         assignAt,
		RevokeAtUpdatged: revokeAtUpdated,
		RevokeAt:         revokeAt,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "timetable to patch wasn't found")
		}
		return nil, status.Error(codes.Internal, "couldn't patch timetable")
	}

	return nil, nil
}

func (s *TimetableService) DeleteTimetable(ctx context.Context, req *ttpb.DeleteTimetableRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("timetableUUID", req.TimetableUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", nil, rbacpb.Verb_DELETE)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for timetable deletion")
	}

	timetableUUID, err := uuid.Parse(req.TimetableUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.timetableRepository.DeleteTimetable(ctx, timetableUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested timetable wasn't found")
		}
		return nil, status.Error(codes.Internal, "couldn't get timetable")
	}

	return nil, nil
}

func (s *TimetableService) ListTimetables(ctx context.Context, req *ttpb.ListTimetablesRequest) (*ttpb.ListTimetablesResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.Int("size", int(req.Size)),
		slog.Int("page", int(req.Page)),
	)
	ctx = logctx.WithLogger(ctx, logger)

	timetables, total, err := s.timetableRepository.ListTimetables(ctx, req.Page, req.Size)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't list timetables")
	}

	timetablesProto := make([]*ttpb.Timetable, len(timetables))
	for i, timetable := range timetables {
		var revokeAt *timestamppb.Timestamp = nil
		if timetable.RevokeAt != nil {
			ra := *timetable.RevokeAt
			revokeAt = timestamppb.New(ra)
		}

		var assignAt *timestamppb.Timestamp = nil
		if timetable.AssignAt != nil {
			aa := *timetable.AssignAt
			revokeAt = timestamppb.New(aa)
		}
		timetablesProto[i] = &ttpb.Timetable{
			UUID:      timetable.UUID.String(),
			Name:      timetable.Name,
			GroupUUID: timetable.GroupUUID.String(),
			RevokeAt:  revokeAt,
			AssignAt:  assignAt,
		}
	}

	return &ttpb.ListTimetablesResponse{
		Page:       req.Page,
		Size:       req.Size,
		Total:      uint64(total),
		Timetables: timetablesProto,
	}, nil
}

func (s *TimetableService) CurrentTimetable(ctx context.Context, req *ttpb.CurrentTimetableRequest) (*ttpb.CurrentTimetableResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("groupUUID", req.GroupUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	groupUUID, err := uuid.Parse(req.GroupUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	timetable, err := s.timetableRepository.CurrentTimetable(ctx, groupUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested timetable wasn't found")
		}
		if errors.Is(err, carefulness.Gone) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Internal, "couldn't get timetable")
	}

	var revokeAt *timestamppb.Timestamp = nil
	if timetable.RevokeAt != nil {
		ra := *timetable.RevokeAt
		revokeAt = timestamppb.New(ra)
	}

	var assignAt *timestamppb.Timestamp = nil
	if timetable.AssignAt != nil {
		aa := *timetable.AssignAt
		revokeAt = timestamppb.New(aa)
	}
	return &ttpb.CurrentTimetableResponse{
		Timetable: &ttpb.Timetable{
			UUID:      timetable.UUID.String(),
			GroupUUID: timetable.GroupUUID.String(),
			Name:      timetable.Name,
			AssignAt:  assignAt,
			RevokeAt:  revokeAt,
		},
	}, nil
}
