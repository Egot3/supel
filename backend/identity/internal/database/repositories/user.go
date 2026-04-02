package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Egot3/supel/backend/identity/internal/database"
	"github.com/Egot3/supel/backend/identity/internal/models"
	passwordutils "github.com/Egot3/supel/backend/identity/internal/passwordUtils"
	"github.com/Egot3/supel/backend/identity/internal/types"
)

func UserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := database.DB.NewSelect().
		Model(&user).
		Where("email = ?", email).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UserById(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	err := database.DB.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func Login(ctx context.Context, email, password string) (string, types.UserRole, error) {
	user, err := UserByEmail(ctx, email)
	if err != nil {
		return "","", err
	}
	if user == nil {
		return "","", errors.New("Invalid credetantials")
	}

	match := passwordutils.CheckPasswordHash(password, user.PasswordHash)
	if !match {
		return "","", errors.New("Invalid credetantials")
	}

	return user.UUID, user.Role, nil
}

func Register(ctx context.Context, email, password string) (string, types.UserRole, error) {
	user, err := UserByEmail(ctx, email)
	if err != nil {
		return "","", err
	}
	if user != nil {
		return "","", errors.New("User with this email alreay exists")
	}

	passwordHash, err := passwordutils.HashPassword(password)
	if err!=nil{
		return "","", err
	}

	user = &models.User{
		Email: email,
		PasswordHash: passwordHash,
	}

	_, err = database.DB.NewInsert().Model(&user).Returning("*").Exec(ctx)

	return user.UUID, user.Role, nil
}


func UpsertUser(ctx context.Context, user models.User) error {
	_, err := database.DB.NewInsert().
		Model(user).
		On("CONFLICT (uuid) DO UPDATE").
		Set("role = ?", user.Role).
		Set("password_hash = ?", user.PasswordHash).
		Set("email = ?", user.Email).
		Exec(ctx)
	return err
}
