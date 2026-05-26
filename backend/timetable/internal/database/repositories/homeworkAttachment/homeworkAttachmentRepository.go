package homeworkattachment

import (
	"context"

	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
)

type HomeworkAttachmentRepository interface {
	HomeworkAttachment(ctx context.Context, homeworkAttachmentUUID uuid.UUID) (*models.HomeworkAttachment, error)
	HomeworkAttachmentsByLessonUUID(ctx context.Context, lessonUUID uuid.UUID) ([]models.HomeworkAttachment, error)
	AttachHomework(ctx context.Context, lessonUUID uuid.UUID, storageKey string) error
	DetachHomework(ctx context.Context, homeworkAttachmentUUID uuid.UUID) error
}
