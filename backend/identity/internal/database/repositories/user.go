package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Egot3/supel/backend/identity/internal/database"
	"github.com/Egot3/supel/backend/identity/internal/models"
)

func User(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User

	err := database.DB.NewSelect().
		Model(user).
		Where("uuid = ?", uuid).
		Scan(ctx, user)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpsertUser(ctx context.Context, user models.User) (error) {
	_, err := database.DB.NewInsert().
		Model(user).
		On("CONFLICT (uuid) DO UPDATE").
		Set("role = ?", user.Role).
		Exec(ctx)
	return err
}