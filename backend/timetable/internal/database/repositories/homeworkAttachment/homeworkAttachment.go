package homeworkattachment

import (
	"context"

	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunHomeworkAttachmentRepository struct {
	db *bun.DB
}

func NewHomeworkAttachmentRepository(i do.Injector) (HomeworkAttachmentRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunHomeworkAttachmentRepository{db: db}, nil
}

func (r *bunHomeworkAttachmentRepository) HomeworkAttachment(ctx context.Context, homeworkAttachmentUUID uuid.UUID) (*models.HomeworkAttachment, error) {
	homeworkAttachment := models.HomeworkAttachment{UUID: homeworkAttachmentUUID}

	err := r.db.NewSelect().Model(&homeworkAttachment).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &homeworkAttachment, nil
}

func (r *bunHomeworkAttachmentRepository) HomeworkAttachmentsByLessonUUID(ctx context.Context, lessonUUID uuid.UUID) ([]models.HomeworkAttachment, error) {
	var homeworkAttachments []models.HomeworkAttachment
	err := r.db.NewSelect().Model(&homeworkAttachments).Where("lesson_uuid = ?", lessonUUID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return homeworkAttachments, nil
}

func (r *bunHomeworkAttachmentRepository) AttachHomework(ctx context.Context, lessonUUID uuid.UUID, storageKey string) error {
	_, err := r.db.NewInsert().Model(&models.HomeworkAttachment{
		StorageKey: storageKey,
		LessonUUID: lessonUUID,
	}).Exec(ctx)

	return err
}

func (r *bunHomeworkAttachmentRepository) DetachHomework(ctx context.Context, homeworkAttachmentUUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.HomeworkAttachment{
		UUID: homeworkAttachmentUUID,
	}).WherePK().Exec(ctx)

	return err
}
