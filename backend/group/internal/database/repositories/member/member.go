package member

import (
	"context"

	"github.com/egot3/supel/backend/group/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunMemberRepository struct {
	db *bun.DB
}

func NewMemberRepository(i do.Injector) (MemberRepository, error) {
	db, err := do.Invoke[*bun.DB](i)
	if err != nil {
		return nil, err
	}

	return &bunMemberRepository{db: db}, nil
}

func (r *bunMemberRepository) AddMember(ctx context.Context, groupUUID, memberUUID uuid.UUID) error {
	_, err := r.db.NewInsert().Model(&models.GroupsMembers{GroupUUID: groupUUID, MemberUUID: memberUUID}).Exec(ctx)
	return err
}

func (r *bunMemberRepository) RemoveMember(ctx context.Context, groupUUID, memberUUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.GroupsMembers{GroupUUID: groupUUID, MemberUUID: memberUUID}).WherePK().Exec(ctx)
	return err
}

func (r *bunMemberRepository) MembersGroups(ctx context.Context, userUUID uuid.UUID, page, size uint32) (uuid.UUIDs, uint64, error) {
	var groups uuid.UUIDs
	total, err := r.db.NewSelect().
		Model((*models.GroupsMembers)(nil)).
		Where("member_uuid = ?", userUUID).
		Limit(int(size)).
		Offset(int(page*size)).
		Column("group_uuid").ScanAndCount(ctx, &groups)
	if err != nil {
		return nil, 0, err
	}

	return groups, uint64(total), nil
}

func (r *bunMemberRepository) ListMembers(ctx context.Context, groupUUID uuid.UUID, page, size uint32) (uuid.UUIDs, uint64, error) {
	var members uuid.UUIDs
	total, err := r.db.NewSelect().
		Model((*models.GroupsMembers)(nil)).
		Where("group_uuid = ?", groupUUID).
		Limit(int(size)).Offset(int(page*size)).
		Column("group_uuid").ScanAndCount(ctx, &members)
	if err != nil {
		return nil, 0, err
	}

	return members, uint64(total), nil
}
