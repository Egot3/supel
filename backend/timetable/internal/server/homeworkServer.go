package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"github.com/egot3/supel/backend/timetable/internal/logctx"
	sanitizationutils "github.com/egot3/supel/backend/timetable/internal/sanitizationUtils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *TimetableService) HomeworkAttachment(ctx context.Context, req *ttpb.HomeworkAttachmentRequest) (*ttpb.HomeworkAttachmentResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("homeworkUUID", req.HomeworkAttachmentUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	attachmentUUID, err := uuid.Parse(req.HomeworkAttachmentUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	attachment, err := s.homeworkAttachmentRepository.HomeworkAttachment(ctx, attachmentUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "requested attachment wasn't found")
		}
		return nil, status.Error(codes.Internal, "couldn't get attachment")
	}

	return &ttpb.HomeworkAttachmentResponse{
		HomeworkAttachment: &ttpb.HomeworkAttachment{
			UUID:       attachment.UUID.String(),
			LessonUUID: attachment.LessonUUID.String(),
			StorageKey: attachment.StorageKey,
			CreatedAt:  timestamppb.New(attachment.CreatedAt),
		},
	}, nil
}

func (s *TimetableService) HomeworkAttachmentByLesson(ctx context.Context, req *ttpb.HomeworkAttachmentByLessonRequest) (*ttpb.HomeworkAttachmentByLessonResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("lessonUUID", req.LessonUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	lessonUUID, err := uuid.Parse(req.LessonUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	attachments, err := s.homeworkAttachmentRepository.HomeworkAttachmentsByLessonUUID(ctx, lessonUUID)

	attachmentsProto := make([]*ttpb.HomeworkAttachment, len(attachments))
	for i, attachment := range attachments {
		attachmentsProto[i] = &ttpb.HomeworkAttachment{
			UUID:       attachment.UUID.String(),
			LessonUUID: attachment.LessonUUID.String(),
			StorageKey: attachment.StorageKey,
			CreatedAt:  timestamppb.New(attachment.CreatedAt),
		}
	}

	return &ttpb.HomeworkAttachmentByLessonResponse{
		HomeworkAttachments: attachmentsProto,
	}, nil
}

func (s *TimetableService) AttachHomework(ctx context.Context, req *ttpb.AttachHomeworkRequest) (*ttpb.AttachHomeworkResponse, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("lessonUUID", req.LessonUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "homework"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for homework attachment")
	}

	lessonUUID, err := uuid.Parse(req.LessonUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	key := sanitizationutils.Slugify(req.Name + req.Mime)
	err = s.homeworkAttachmentRepository.AttachHomework(ctx, lessonUUID, key)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't create homework attachment")
	}

	URL, err := s.storage.PUTurl(ctx, key, req.Mime)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't generate url for attachment")
	}

	return &ttpb.AttachHomeworkResponse{
		UploadUrl: URL,
		FileKey:   key,
	}, nil
}

func (s *TimetableService) DetachHomework(ctx context.Context, req *ttpb.DetachHomeworkRequest) (*emptypb.Empty, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "handler"),
		slog.String("homeworkUUID", req.HomeworkAttachmentUUID),
	)
	ctx = logctx.WithLogger(ctx, logger)

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "unable to get own uuid from token")
	}

	subscope := "homework"
	can, err := s.grpcClient.HasPermission(ctx, ownUUID, "timetable", &subscope, rbacpb.Verb_DELETE)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't check user's permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "user doesn't have enough permissions for homework detachment")
	}

	attachmentUUID, err := uuid.Parse(req.HomeworkAttachmentUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad uuid")
	}

	attachment, err := s.homeworkAttachmentRepository.HomeworkAttachment(ctx, attachmentUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "couldn't find homework")
		}
		return nil, status.Error(codes.Internal, "couldn't get homework")
	}

	err = s.storage.DeleteFromKey(ctx, attachment.StorageKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't delete homework")
	}

	err = s.homeworkAttachmentRepository.DetachHomework(ctx, attachmentUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "couldn't find homework")
		}
		return nil, status.Error(codes.Internal, "couldn't detach homework")
	}

	return nil, nil
}
