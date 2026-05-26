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

func (s *TimetableService) Subject(ctx context.Context, req *ttpb.SubjectRequest) (*ttpb.SubjectResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("subjectUUID", req.SubjectUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	subjectUUID, err := uuid.Parse(req.SubjectUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	subject, err := s.subjectRepository.Subject(ctx, subjectUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested subject wasn't found")
		}
		if errors.Is(err, carefulness.Gone) {
			metadata.AppendToOutgoingContext(ctx, "reponse-code", "501")
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, status.Error(codes.Internal, "couldn't get subject")
	}

	return &ttpb.SubjectResponse{
		Subject: &ttpb.Subject{
			Name: subject.Name,
			UUID: subject.UUID.String(),
		},
	}, nil
}

func (s *TimetableService) CreateSubject(ctx context.Context, req *ttpb.CreateSubjectRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("name", req.Name),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "subject"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for subject creation")
	}

	err = s.subjectRepository.CreateSubject(ctx, req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't create subject")
	}

	return nil, nil
}

func (s *TimetableService) PatchSubject(ctx context.Context, req *ttpb.PatchSubjectRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("subjectUUID", req.SubjectUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "subject"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_PATCH)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for subject patching")
	}

	subjectUUID, err := uuid.Parse(req.SubjectUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.subjectRepository.PatchSubject(ctx, models.SubjectPatched{
		UUID: subjectUUID,
		Name: req.Name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "couldn't find subject to patch")
		}
		return nil, status.Error(codes.Internal, "couldn't patch subject")
	}

	return nil, nil
}

func (s *TimetableService) DeleteSubject(ctx context.Context, req *ttpb.DeleteSubjectRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("subjectUUID", req.SubjectUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "subject"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_DELETE)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for subject deletion")
	}

	subjectUUID, err := uuid.Parse(req.SubjectUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	err = s.subjectRepository.DeleteSubject(ctx, subjectUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "couldn't find subject to delete")
		}
		return nil, status.Error(codes.Internal, "couldn't delete subject")
	}

	return nil, nil
}

func (s *TimetableService) ListSubjects(ctx context.Context, req *ttpb.ListSubjectsRequest) (*ttpb.ListSubjectsResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.Int("size", int(req.Size)),
		slog.Int("page", int(req.Page)),
	)
	ctx = logctx.WithLogger(ctx, logger)

	subjects, total, err := s.subjectRepository.ListSubjects(ctx, req.Page, req.Size)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't list subjects")
	}

	subjectsProto := make([]*ttpb.Subject, len(subjects))
	for i, subject := range subjects {
		subjectsProto[i] = &ttpb.Subject{
			UUID: subject.UUID.String(),
			Name: subject.Name,
		}
	}

	return &ttpb.ListSubjectsResponse{
		Page:     req.Page,
		Size:     req.Size,
		Total:    uint64(total),
		Subjects: subjectsProto,
	}, nil
}
