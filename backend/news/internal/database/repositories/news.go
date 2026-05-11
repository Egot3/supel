package repositories

import (
	"context"
	"database/sql"

	"github.com/Egot3/supel/backend/news/internal/database"
	"github.com/Egot3/supel/backend/news/internal/models"
	"github.com/uptrace/bun"
)

func CreateNew(ctx context.Context, new models.New, images []string) (*models.New, error) {

	err := database.DB.RunInTx(ctx, &sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.NewInsert().Model(&new).Returning("*").Exec(ctx)
			if err != nil {
				return err
			}

			if len(images) > 0 {
				for i, imageKey := range images {
					_, err := tx.NewInsert().Model(&models.NewsImages{
						NewUUID:  new.NewUUID,
						FileKey:  imageKey,
						Position: i,
					}).Exec(ctx)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	return &new, err
}

func NewImagesByUUId(ctx context.Context, uuid string) (imageKeys []string, err error) {
	err = database.DB.NewSelect().
		Model((*models.NewsImages)(nil)).
		Column("file_key").
		Where("new_uuid = ?", uuid).
		OrderBy("position", bun.OrderAsc).
		Scan(ctx, &imageKeys)
	if err != nil {
		return nil, err
	}

	return imageKeys, nil
}

func NewByUUID(ctx context.Context, uuid string) (*models.New, error) {
	var New *models.New
	_, err := database.DB.NewSelect().
		Model(&New).
		Where("new_uuid = ?", uuid).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return New, nil
}

func NewBulk(ctx context.Context, page, size int) ([]models.New, int, error) {
	var found []models.New
	total, err := database.DB.NewSelect().
		Model(&found).
		Limit(size).
		Offset(page * size).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return found, total, err
}

func NewBulkByUser(ctx context.Context, page, size int, userUUID string) ([]models.New, int, error) {
	var found []models.New
	total, err := database.DB.NewSelect().
		Model(&found).
		Limit(size).
		Where("user_uuid = ?", userUUID).
		Offset(page * size).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return found, total, err
}

func IsCreator(ctx context.Context, userUUID, newUUID string) (is bool, err error) {
	return database.DB.NewSelect().
		Model((*models.New)(nil)).
		Where("new_uuid = ?", newUUID).
		Where("user_uuid = ?", userUUID).
		Exists(ctx)
}

func DeleteNew(ctx context.Context, newUUID string) error {
	_, err := database.DB.NewDelete().
		Model((*models.New)(nil)).
		Where("new_uuid = ?", newUUID).
		Exec(ctx)
	return err
}
