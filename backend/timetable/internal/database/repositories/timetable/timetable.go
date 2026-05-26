package timetable

import (
	"context"
	"time"

	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunTimetableRepository struct {
	db *bun.DB
}

func NewTimetableRepository(i do.Injector) (TimetableRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunTimetableRepository{db: db}, nil
}

func (r *bunTimetableRepository) Timetable(ctx context.Context, timetableUUID uuid.UUID) (*models.Timetable, error) {
	timetable := models.Timetable{UUID: timetableUUID}

	err := r.db.NewSelect().Model(&timetable).WherePK().Relation("TimetableEntry").WhereAllWithDeleted().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &timetable, nil
}

func (r *bunTimetableRepository) CreateTimetable(ctx context.Context, groupUUID uuid.UUID, name string, assignAt, revokeAt *time.Time) error {
	_, err := r.db.NewInsert().
		Model(&models.Timetable{Name: name, GroupUUID: groupUUID, AssignAt: assignAt, RevokeAt: revokeAt}).
		Exec(ctx)
	return err
}

func (r *bunTimetableRepository) PatchTimetable(ctx context.Context, patched models.TimetablePatched) error {
	isUpdated := false
	query := r.db.NewUpdate().Model(&models.Timetable{UUID: patched.UUID}).WherePK()

	if patched.GroupUUID != uuid.Nil {
		query = query.Set("group_uuid = ?", patched.GroupUUID)
		isUpdated = true
	}
	if patched.Name != nil {
		query = query.Set("name = ?", patched.Name)
		isUpdated = true
	}
	if patched.AssignAtUpdated {
		query = query.Set("assign_at = ?", patched.AssignAt)
		isUpdated = true
	}
	if patched.RevokeAtUpdatged {
		query = query.Set("revoke_at = ?", patched.RevokeAt)
		isUpdated = true
	}

	if !isUpdated {
		return nil
	}

	_, err := query.Exec(ctx)
	return err
}

func (r *bunTimetableRepository) DeleteTimetable(ctx context.Context, timetableUUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.Timetable{UUID: timetableUUID}).WherePK().Exec(ctx)
	return err
}

func (r *bunTimetableRepository) ListTimetables(ctx context.Context, page, size uint32) ([]models.Timetable, int, error) {
	var timetables []models.Timetable
	total, err := r.db.NewSelect().Model(&timetables).Limit(int(size)).Offset(int(page * size)).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return timetables, total, nil
}

func (r *bunTimetableRepository) CurrentTimetable(ctx context.Context, groupUUID uuid.UUID) (*models.Timetable, error) {
	var timetable models.Timetable

	err := r.db.NewSelect().Model(&timetable).
		Where("group_uuid = ?", groupUUID).WhereAllWithDeleted().Where("? BETWEEN assign_at AND revoke_at", time.Now()).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &timetable, nil
}

func (r *bunTimetableRepository) TimetableByDate(ctx context.Context, groupUUID uuid.UUID, date time.Time) (*models.Timetable, error) {
	var timetable models.Timetable

	err := r.db.NewSelect().Model(&timetable).
		Where("group_uuid = ?", groupUUID).WhereAllWithDeleted().Where("? BETWEEN assign_at AND revoke_at", date).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &timetable, nil
}
