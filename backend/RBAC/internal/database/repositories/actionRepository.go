package repositories

import (
	"context"

	"github.com/google/uuid"
)

type ActionRepository interface {
	AddActions(ctx context.Context, actionUUIDs []uuid.UUID, roleUUID uuid.UUID) error
	RevokeActions(ctx context.Context, actionUUIDs []uuid.UUID, roleUUID uuid.UUID) error
}
