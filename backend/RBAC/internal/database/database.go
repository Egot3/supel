package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func InitDB(i do.Injector) (*bun.DB, error) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER_NAME")) //лишь бы не запутаться

	// log.Printf("dsn: %v\n", dsn)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	DB := bun.NewDB(sqldb, pgdialect.New())

	time.Sleep(2 * time.Second)
	for i := range 5 {
		if err := DB.Ping(); err != nil {
			log.Printf("Попытка %d: Пингов не прошло: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	} //почему бы и нет

	if err := DB.Ping(); err != nil {
		log.Printf("\nтут уже можно не продолжать, база легла\n")
		return nil, err
	}

	sqldb.SetMaxOpenConns(50)
	sqldb.SetMaxIdleConns(20)

	log.Printf("ДБ стоит")
	return DB, nil
}
