package repositories

import (
	"context"

	"github.com/Egot3/supel/backend/user/internal/database"
	"github.com/Egot3/supel/backend/user/internal/models"
)

func CreateUser(ctx context.Context, nickname, uuid string, avatarKey string) error {
	_, err := database.DB.NewInsert().Model(&models.User{Nickname: nickname, UUID: uuid, AvatarKey: avatarKey}).Exec(ctx)
	return err
}

func DeleteUser(ctx context.Context, uuid string) error {
	_, err := database.DB.NewDelete().Model(&models.User{UUID: uuid}).WherePK().Exec(ctx)
	return err
}

func GetUser(ctx context.Context, uuid string) (*models.User, error) {
	user := &models.User{
		UUID: uuid,
	}

	err := database.DB.NewSelect().Model(user).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func PatchUser(ctx context.Context, patched *models.UpdateUser) error {
	updateQuery := database.DB.NewUpdate().Table("users")
	if patched.Nickname != nil {
		updateQuery = updateQuery.Set("nickname = ?", patched.Nickname)
	}
	if patched.Description != nil {
		updateQuery = updateQuery.Set("description = ?", patched.Description)
	}
	if patched.Status != nil {
		updateQuery = updateQuery.Set("status = ?", patched.Status)
	}
	if patched.StatusExpiration != nil {
		updateQuery = updateQuery.Set("status_expiration = ?", patched.StatusExpiration)
	}
	if patched.StatusReactionKey != nil {
		updateQuery = updateQuery.Set("status_reaction_key = ?", patched.StatusReactionKey)
	}

	_, err := updateQuery.Where("uuid = ?", patched.UUID).
		Exec(ctx)
	return err
}

func GetAvatarKey(ctx context.Context, uuid string) (string, error) {
	var key string
	err := database.DB.NewSelect().
		Model(&models.User{UUID: uuid}).
		WherePK().
		Column("avatar_key").
		Scan(ctx, &key)

	return key, err
}
