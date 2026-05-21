package curator

import (
	"context"

	"github.com/google/uuid"
)

type CuratorRepository interface {
	AssignCuratorToSenior(ctx context.Context, seniorUUID, subordinateUUID uuid.UUID) error
	RevokeCuratorFromSenior(ctx context.Context, seniorUUID, subordinateUUID uuid.UUID) error
	AssignCuratorToGroup(ctx context.Context, curatorUUID, groupUUID uuid.UUID) error
	RevokeCuratorFromGroup(ctx context.Context, curatorUUID, groupUUID uuid.UUID) error
	AddCurator(ctx context.Context, seniorUUID, curatorUUID, groupUUID uuid.UUID) error
	CanEdit(ctx context.Context, curatorUUID, groupUUID uuid.UUID) (bool, error)
	WillCycle(ctx context.Context, seniorUUID, subordinateUUID uuid.UUID) (bool, error)
	RevokeCurator(ctx context.Context, curatorUUID uuid.UUID) error
}
