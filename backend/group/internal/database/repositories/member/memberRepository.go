package member

import (
	"context"

	"github.com/google/uuid"
)

type MemberRepository interface {
	AddMember(ctx context.Context, groupUUID, memberUUID uuid.UUID) error
	RemoveMember(ctx context.Context, groupUUID, memberUUID uuid.UUID) error
	MembersGroups(ctx context.Context, userUUID uuid.UUID, page, size int32) (uuid.UUIDs, int64, error)
	ListMembers(ctx context.Context, groupUUID uuid.UUID, page, size int32) (uuid.UUIDs, int64, error)
}
