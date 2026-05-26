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

func (s *TimetableService) Lesson(ctx context.Context, req *ttpb.LessonRequest) (*ttpb.LessonResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("lessonUUID", req.LessonUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	lessonUUID, err := uuid.Parse(req.LessonUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	lesson, err := s.lessonRepository.Lesson(ctx, lessonUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested lesson wasn't found")
		}
		if errors.Is(err, carefulness.Gone) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Internal, "couldn't get timetable")
	}

	return &ttpb.LessonResponse{
		Lesson: &ttpb.Lesson{
			UUID:               lesson.UUID.String(),
			TimetableEntryUUID: lesson.TimetableEntryUUID.String(),
			Date:               timestamppb.New(lesson.Date),
			Cancelled:          lesson.Cancelled,
			CreatedAt:          timestamppb.New(lesson.CreatedAt),
		},
	}, nil
}

func (s *TimetableService) CreateLesson(ctx context.Context, req *ttpb.CreateLessonRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("timetableEntryUUID", req.TimetableEntryUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "lesson"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for lesson creation")
	}

	timetableEntryUUID, err := uuid.Parse(req.TimetableEntryUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.lessonRepository.CreateLesson(ctx, timetableEntryUUID, req.Date.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't create Lesson")
	}

	return nil, nil
}

func (s *TimetableService) PatchLesson(ctx context.Context, req *ttpb.PatchLessonRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("lessonUUID", req.LessonUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "lesson"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_PATCH)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for lesson patching")
	}

	lessonUUID, err := uuid.Parse(req.LessonUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	var date *time.Time = nil
	if req.Date != nil {
		d := req.Date.AsTime()
		date = &d
	}

	var tableEntryUUID uuid.UUID = uuid.Nil
	if req.TimetableEntryUUID != nil {
		tableEntryUUID, err = uuid.Parse(*req.TimetableEntryUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Bad uuid")
		}
	}

	err = s.lessonRepository.PatchLesson(ctx, models.LessonPatched{
		UUID:               lessonUUID,
		Date:               date,
		Cancelled:          req.Cancelled,
		TimetableEntryUUID: tableEntryUUID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "given lesson wasn't found")
		}
		return nil, status.Error(codes.Internal, "couldn't patch lesson")
	}

	return nil, nil
}

func (s *TimetableService) DeleteLesson(ctx context.Context, req *ttpb.DeleteLessonRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("lessonUUID", req.LessonUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	lessonUUID, err := uuid.Parse(req.LessonUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.lessonRepository.DeleteLesson(ctx, lessonUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested lesson wasn't found")
		}
		if errors.Is(err, carefulness.Gone) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Internal, "couldn't delete lessen")
	}

	return nil, nil
}

func (s *TimetableService) ListLessons(ctx context.Context, req *ttpb.ListLessonsRequest) (*ttpb.ListLessonsResponse, error) {
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
		return nil, status.Error(codes.Internal, "couldn't get timetable for requested lessons")
	}

	lessons := make([]*ttpb.Lesson, len(timetable.TimetableEntries))
	for i, timetableEntry := range timetable.TimetableEntries {
		lesson, err := s.lessonRepository.LessonByEntry(ctx, timetableEntry.UUID)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't get lesson for entry")
		}
		lessons[i] = &ttpb.Lesson{
			UUID:               lesson.UUID.String(),
			TimetableEntryUUID: lesson.TimetableEntryUUID.String(),
			Date:               timestamppb.New(lesson.Date),
			Cancelled:          lesson.Cancelled,
			CreatedAt:          timestamppb.New(lesson.CreatedAt),
		}
	}

	return &ttpb.ListLessonsResponse{
		Lessons: lessons,
	}, nil
}
