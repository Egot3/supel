package repositories

import (
	"context"
	"time"

	"github.com/Egot3/supel/backend/identity/internal/database"
	"github.com/Egot3/supel/backend/identity/internal/models"
)

func UpdateSyncing(ctx context.Context, newSync models.Sync) (error) {
	exists, err := database.DB.NewSelect().Model(newSync).WherePK().Exists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		if _, err := database.DB.NewInsert().Model(newSync).Exec(ctx); err != nil {
			return err
		}
	}

	_, err = database.DB.NewUpdate().Model(newSync).Set("last_synced_at = ?", time.Now()).WherePK().Exec(ctx)

	return err
}