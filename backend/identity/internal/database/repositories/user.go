package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	carefulness "github.com/Egot3/supel/backend/identity/internal"
	"github.com/Egot3/supel/backend/identity/internal/database"
	"github.com/Egot3/supel/backend/identity/internal/models"
	passwordutils "github.com/Egot3/supel/backend/identity/internal/passwordUtils"
	"github.com/uptrace/bun"
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
		Where("uuid = ?", id).
		Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func Login(ctx context.Context, email, password string) (string, error) {
	user, err := UserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", carefulness.InvalidCreditantials
	}

	match := passwordutils.CheckPasswordHash(password, user.PasswordHash)
	if !match {
		return "", carefulness.InvalidCreditantials
	}

	return user.UUID, nil
}

func Register(ctx context.Context, email, password string) (string, error) {
	log.Printf("registering user with email: %v", email)
	user, err := UserByEmail(ctx, email)
	if err != nil {
		log.Printf("error while fetching user: %v", err)
		return "", fmt.Errorf("error while fetching user")
	}
	if user != nil {
		log.Printf("user was found")
		return "", carefulness.ErrEmailAlreadyExists
	}

	passwordHash, err := passwordutils.HashPassword(password)
	if err != nil {
		log.Printf("hashing password failed: %v", err)
		return "", err
	}

	user = &models.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	_, err = database.DB.NewInsert().Model(user).Returning("*").Exec(ctx)
	if err != nil {
		log.Println("error while inserting user: ", err.Error())
		return "", err
	}

	return user.UUID, nil
}

func UpsertUser(ctx context.Context, user models.User) error {
	_, err := database.DB.NewInsert().
		Model(user).
		On("CONFLICT (uuid) DO UPDATE").
		Set("password_hash = ?", user.PasswordHash).
		Set("email = ?", user.Email).
		Exec(ctx)
	return err
}

func DisableUserTx(ctx context.Context, tx bun.Tx, uuid string) error {
	resp, err := tx.NewUpdate().
		Model((*models.User)(nil)).
		Where("uuid = ?", uuid).
		Set("is_active = false").Exec(ctx)
	if err != nil {
		return err
	}

	if val, _ := resp.RowsAffected(); val == 0 {
		return carefulness.UserNotFound
	}

	return nil
}
