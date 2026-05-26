package timetableentry

import (
	"context"

	"github.com/egot3/supel/backend/timetable/internal/carefulness"
	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunTimetableEntryRepository struct {
	db *bun.DB
}

func NewTimetableEntryRepository(i do.Injector) (TimetableEntryRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunTimetableEntryRepository{db: db}, nil
}

func (r *bunTimetableEntryRepository) TimetableEntry(ctx context.Context, timetableEntryUUID uuid.UUID) (*models.TimetableEntry, error) {
	timetableEntry := models.TimetableEntry{UUID: timetableEntryUUID}
	err := r.db.NewSelect().Model(&timetableEntry).WherePK().WhereAllWithDeleted().Scan(ctx)
	if err != nil {
		return nil, err
	}

	if timetableEntry.DeletedAt != nil {
		return nil, carefulness.Gone
	}

	return &timetableEntry, nil
}

func (r *bunTimetableEntryRepository) TimetableEntriesByTimetable(ctx context.Context, timetableUUID uuid.UUID) ([]models.TimetableEntry, error) {
	var timetableEntries []models.TimetableEntry
	err := r.db.NewSelect().Model(&timetableEntries).Where("timetable_uuid = ?", timetableUUID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return timetableEntries, nil
}

func (r *bunTimetableEntryRepository) CreateTimetableEntry(ctx context.Context, timetableUUID, periodUUID, subjectUUID, teacherUUID uuid.UUID, dayOfWeek int16, place string) error {
	_, err := r.db.NewInsert().Model(&models.TimetableEntry{
		TimetableUUID: timetableUUID,
		PeriodUUID:    periodUUID,
		SubjectUUID:   subjectUUID,
		TeacherUUID:   teacherUUID,
		DayOfWeek:     dayOfWeek,
		Place:         place,
	}).Exec(ctx)
	return err
}

func (r *bunTimetableEntryRepository) PatchTimetableEntry(ctx context.Context, patched models.TimetableEntryPatched) error {
	isUpdated := false

	query := r.db.NewUpdate().Model(&models.TimetableEntry{UUID: patched.UUID}).WherePK()
	if patched.TimetableUUID != uuid.Nil {
		query = query.Set("timetable_uuid = ?", patched.TimetableUUID)
		isUpdated = true
	}
	if patched.PeriodUUID != uuid.Nil {
		query = query.Set("period_uuid = ?", patched.PeriodUUID)
		isUpdated = true
	}
	if patched.SubjectUUID != uuid.Nil {
		query = query.Set("subject_uuid = ?", patched.SubjectUUID)
		isUpdated = true
	}
	if patched.TeacherUpdated {
		query = query.Set("teacher_uuid = ?", patched.TeacherUUID)
		isUpdated = true

	}
	if patched.DayOfWeek != 0 {
		query = query.Set("day_of_week = ?", patched.DayOfWeek)
		isUpdated = true

	}
	if patched.Place != nil {
		query = query.Set("place = ?", patched.Place)
		isUpdated = true

	}

	if !isUpdated {
		return nil
	}

	_, err := query.Exec(ctx)
	return err
}

func (r *bunTimetableEntryRepository) DeleteTimetableEntry(ctx context.Context, timetableEntryUUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.TimetableEntry{UUID: timetableEntryUUID}).WherePK().Exec(ctx)
	return err
}
