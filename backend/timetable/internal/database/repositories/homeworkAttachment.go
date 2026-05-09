package repositories

import (
	"context"

	"github.com/Egot3/supel/backend/timetable/internal/models"
	"github.com/uptrace/bun"
)

type bunHomeworkAttachmentRepository struct {
	db *bun.DB
}

func NewHomeworkAttachmentRepository(db *bun.DB) HomeworkAttachmentRepository {
	return &bunHomeworkAttachmentRepository{db: db}
}

func (r *bunHomeworkAttachmentRepository) CreateHomeworkAttachment(ctx context.Context, hwa models.HomeworkAttachment) (*models.HomeworkAttachment, error) {
	_, err := r.db.NewInsert().Model(&hwa).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &hwa, err
}

func (r *bunHomeworkAttachmentRepository) DeleteHomeworkAttachment(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().Model(&models.HomeworkAttachment{FileUUID: uuid}).WherePK().Exec(ctx)
	return err
}

func (r *bunHomeworkAttachmentRepository) HomeworkAttachmentByUUID(ctx context.Context, uuid string) (*models.HomeworkAttachment, error) {
	hwa := models.HomeworkAttachment{FileUUID: uuid}
	err := r.db.NewSelect().Model(&hwa).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &hwa, err
}

func (r *bunHomeworkAttachmentRepository) HomeworkAttachmentKeyByUUID(ctx context.Context, uuid string) (string, error) {
	var key string
	err := r.db.NewSelect().Model(models.HomeworkAttachment{FileUUID: uuid}).WherePK().Column("storage_key").Scan(ctx, key)

	if err != nil {
		return "", err
	}

	return key, nil
}

func (r *bunHomeworkAttachmentRepository) HomeworkAttachmentKeysByConcreteUUID(ctx context.Context, uuid string) ([]string, error) {
	var keys []string
	err := r.db.NewSelect().
		Model((*models.HomeworkAttachment)(nil)).
		Where("concrete_uuid = ?", uuid).
		Column("storage_key").Scan(ctx, &keys)

	if err != nil {
		return nil, err
	}

	return keys, nil
}
