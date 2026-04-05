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
