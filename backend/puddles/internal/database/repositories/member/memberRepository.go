package member

import (
	"context"

	"github.com/egot3/supel/backend/puddles/internal/models"
	"github.com/google/uuid"
)

type MemberRepository interface {
	AddMember(ctx context.Context, memberUUID uuid.UUID, adderUUID *uuid.UUID, puddleUUID uuid.UUID) error
	ListMembersPuddles(ctx context.Context, memberUUID uuid.UUID, page, size uint32) (uuid.UUIDs, int, error)
	ListMembersPuddlesIntersections(ctx context.Context, memberAUUID, memberBUUID uuid.UUID, page, size uint32) (uuid.UUIDs, int, error)
	ListAddersAddors(ctx context.Context, adderUUID uuid.UUID, page, size uint32) (uuid.UUIDs, int, error)
	PuddleMember(ctx context.Context, puddleUUID, memberUUID uuid.UUID) (*models.PuddlesMembers, error)
}
