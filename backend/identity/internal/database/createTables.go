package database

import (
	"context"
	"log"

	"github.com/Egot3/supel/backend/identity/internal/models"
)

func CreateTables() {
	ctx := context.Background()

	_, err := DB.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)

	if err != nil {
		log.Printf("|ПЛОХИЕ НОВОСТИ| таблица юзеров не была создана: %v\n", err)
	}

}
