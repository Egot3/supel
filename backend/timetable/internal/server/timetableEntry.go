package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

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
)

func (s *TimetableService) TimetableEntry(ctx context.Context, req *ttpb.TimetableEntryRequest) (*ttpb.TimetableEntryResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("entryUUID", req.TimetableEntryUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	entryUUID, err := uuid.Parse(req.TimetableEntryUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	entry, err := s.timetableEntryRepository.TimetableEntry(ctx, entryUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested entry wasn't found")
		}
		if errors.Is(err, carefulness.Gone) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Internal, "couldn't get entry")
	}

	var teacherUUID *string = nil
	if entry.TeacherUUID != uuid.Nil {
		tus := entry.TeacherUUID.String()
		teacherUUID = &tus
	}
	return &ttpb.TimetableEntryResponse{
		TimetableEntry: &ttpb.TimetableEntry{
			UUID:          entry.UUID.String(),
			TimetableUUID: entry.TimetableUUID.String(),
			PeriodUUID:    entry.PeriodUUID.String(),
			SubjectUUID:   entry.SubjectUUID.String(),
			TeacherUUID:   teacherUUID,
			Place:         entry.Place,
			DayOfWeek:     ttpb.DayOfWeek(entry.DayOfWeek),
		},
	}, nil
}

func (s *TimetableService) CreateTimetableEntry(ctx context.Context, req *ttpb.CreateTimetableEntryRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("subjectUUID", req.SubjectUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "entry"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for entry creation")
	}

	timetableUUID, err := uuid.Parse(req.TimetableUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	var teacherUUID uuid.UUID = uuid.Nil
	if req.TeacherUUID != nil {
		teacherUUID, err = uuid.Parse(*req.TeacherUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Bad uuid")
		}
	}

	periodUUID, err := uuid.Parse(req.PeriodUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}
	subjectUUID, err := uuid.Parse(req.SubjectUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.timetableEntryRepository.CreateTimetableEntry(ctx, timetableUUID, periodUUID, subjectUUID, teacherUUID, int16(req.DayOfWeek), req.Place)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't create entry")
	}

	return nil, nil
}

func (s *TimetableService) PatchTimetableEntry(ctx context.Context, req *ttpb.PatchTimetableEntryRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("entryUUID", req.TimetableEntryUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "entry"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_PATCH)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for lesson patching")
	}

	entryUUID, err := uuid.Parse(req.TimetableEntryUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	teacherUpdated := false
	var teacherUUID uuid.UUID = uuid.Nil
	if req.TeacherUUID != nil {
		teacherUUID, err = uuid.Parse(*req.TeacherUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Bad uuid")
		}
		teacherUpdated = true
	}
	var subjectUUID uuid.UUID = uuid.Nil
	if req.SubjectUUID != nil {
		subjectUUID, err = uuid.Parse(*req.SubjectUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Bad uuid")
		}
	}
	var tableUUID uuid.UUID = uuid.Nil
	if req.TimetableUUID != nil {
		tableUUID, err = uuid.Parse(*req.TimetableUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Bad uuid")
		}
	}
	var periodUUID uuid.UUID = uuid.Nil
	if req.PeriodUUID != nil {
		periodUUID, err = uuid.Parse(*req.PeriodUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Bad uuid")
		}
	}

	err = s.timetableEntryRepository.PatchTimetableEntry(ctx, models.TimetableEntryPatched{
		UUID:           entryUUID,
		TimetableUUID:  tableUUID,
		PeriodUUID:     periodUUID,
		TeacherUpdated: teacherUpdated,
		TeacherUUID:    teacherUUID,
		SubjectUUID:    subjectUUID,
		Place:          req.Place,
		DayOfWeek:      int16(req.DayOfWeek),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested to patch entry wasn't found")
		}
		return nil, status.Error(codes.Internal, "couldn't patch entry")
	}

	return nil, nil
}

func (s *TimetableService) DeleteTimetableEntry(ctx context.Context, req *ttpb.DeleteTimetableEntryRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("entryUUID", req.TimetableEntryUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "entry"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for entry creation")
	}

	entryUUID, err := uuid.Parse(req.TimetableEntryUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.timetableEntryRepository.DeleteTimetableEntry(ctx, entryUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested to delete entry wasn't found")
		}
		return nil, status.Error(codes.Internal, "couldn't delete an entry")
	}

	return nil, nil
}

func (s *TimetableService) TimetableEntriesByTimetable(ctx context.Context, req *ttpb.TimetableEntriesByTimetableRequest) (*ttpb.TimetableEntriesByTimetableResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("tableUUID", req.TimetableUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	tableUUID, err := uuid.Parse(req.TimetableUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	table, err := s.timetableRepository.Timetable(ctx, tableUUID)
	entries := make([]*ttpb.TimetableEntry, len(table.TimetableEntries))

	for i, entry := range table.TimetableEntries {

		var teacherUUID *string = nil
		if entry.TeacherUUID != uuid.Nil {
			tus := entry.TeacherUUID.String()
			teacherUUID = &tus
		}
		entries[i] = &ttpb.TimetableEntry{
			UUID:          entry.UUID.String(),
			TimetableUUID: entry.TimetableUUID.String(),
			PeriodUUID:    entry.PeriodUUID.String(),
			SubjectUUID:   entry.SubjectUUID.String(),
			TeacherUUID:   teacherUUID,
			Place:         entry.Place,
			DayOfWeek:     ttpb.DayOfWeek(entry.DayOfWeek),
		}
	}

	return &ttpb.TimetableEntriesByTimetableResponse{
		TimetableEntries: entries,
	}, nil
}
