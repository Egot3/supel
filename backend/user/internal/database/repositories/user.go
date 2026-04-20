package repositories

import (
	"context"

	"github.com/Egot3/supel/backend/user/internal/database"
	"github.com/Egot3/supel/backend/user/internal/models"
)

func CreateUser(ctx context.Context, nickname, uuid string) error {
	_, err := database.DB.NewInsert().Model(models.User{Nickname: nickname, UUID: uuid}).Exec(ctx)
	return err
}

func DeleteUser(ctx context.Context, uuid string) error {
	_, err := database.DB.NewDelete().Model(models.User{UUID: uuid}).WherePK().Exec(ctx)
	return err
}

func GetUser(ctx context.Context, uuid string) (*models.User, error) {
	user := &models.User{
		UUID: uuid,
	}

	err := database.DB.NewSelect().Model(&user).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func PatchUser(ctx context.Context, patched *models.User) error {
	_, err := database.DB.NewUpdate().Model(patched).WherePK().Exec(ctx)
	return err
}
