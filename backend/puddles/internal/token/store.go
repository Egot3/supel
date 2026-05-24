package token

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do/v2"
)

const (
	TICKETTTL = 30 * time.Second
)

type redisStore struct {
	rdb *redis.Client
}

type Store interface {
	IssueWSTicket(ctx context.Context, userUUID uuid.UUID) (string, error)
	ConsumeWSTicket(ctx context.Context, ticket string) (uuid.UUID, error)
}

func NewStore(i do.Injector) Store {
	rdb := do.MustInvoke[*redis.Client](i)
	return &redisStore{rdb: rdb}
}

func generateToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(b[:]), nil
}

func (s *redisStore) IssueWSTicket(ctx context.Context, userUUID uuid.UUID) (string, error) {
	ticket, err := generateToken()
	if err != nil {
		return "", err
	}

	if err := s.rdb.Set(ctx, fmt.Sprintf("chat:ws:ticket:%v", ticket), userUUID, TICKETTTL).Err(); err != nil {
		return "", err
	}

	return ticket, nil
}

func (s *redisStore) ConsumeWSTicket(ctx context.Context, ticket string) (uuid.UUID, error) {
	userUUID, err := s.rdb.GetDel(ctx, fmt.Sprintf("chat:ws:ticket:%v", ticket)).Result()
	if err == redis.Nil {
		return uuid.Nil, nil
	}
	if err != nil {
		return uuid.Nil, err
	}
	parsedUserUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return uuid.Nil, err
	}

	return parsedUserUUID, nil
}
