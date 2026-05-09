package repositories

import (
	"context"

	"github.com/Egot3/supel/backend/timetable/internal/models"
)

type HomeworkAttachmentRepository interface {
	CreateHomeworkAttachment(ctx context.Context, hwa models.HomeworkAttachment) (*models.HomeworkAttachment, error)
	DeleteHomeworkAttachment(ctx context.Context, uuid string) error
	HomeworkAttachmentByUUID(ctx context.Context, uuid string) (*models.HomeworkAttachment, error)
	HomeworkAttachmentKeyByUUID(ctx context.Context, uuid string) (string, error)
	HomeworkAttachmentKeysByConcreteUUID(ctx context.Context, uuid string) ([]string, error)
}
