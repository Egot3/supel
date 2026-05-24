package redisdatabase

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

func New(i do.Injector) (*redis.Client, error) {
	cfg := do.MustInvoke[Config](i)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,

		PoolSize:        10,
		MinIdleConns:    1,
		ConnMaxIdleTime: 2 * time.Minute,

		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
