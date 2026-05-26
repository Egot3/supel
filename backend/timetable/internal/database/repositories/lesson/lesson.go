package lesson

import (
	"context"
	"time"

	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunLessonRepository struct {
	db *bun.DB
}

func NewLessonRepository(i do.Injector) (LessonRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunLessonRepository{db: db}, nil
}

func (r *bunLessonRepository) Lesson(ctx context.Context, lessonUUID uuid.UUID) (*models.Lesson, error) {
	lesson := models.Lesson{UUID: lessonUUID}
	err := r.db.NewSelect().Model(&lesson).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &lesson, nil
}

func (r *bunLessonRepository) CreateLesson(ctx context.Context, timetableEntryUUID uuid.UUID, date time.Time) error {
	_, err := r.db.NewInsert().Model(&models.Lesson{UUID: timetableEntryUUID, Date: date, Cancelled: false}).Exec(ctx)
	return err
}

func (r *bunLessonRepository) PatchLesson(ctx context.Context, pathced models.LessonPatched) error {
	isUpdated := false

	query := r.db.NewUpdate().Model(&models.Lesson{UUID: pathced.UUID, Cancelled: pathced.Cancelled}).WherePK()
	if pathced.TimetableEntryUUID != uuid.Nil {
		query = query.Set("timetable_entry_uuid = ?", pathced.TimetableEntryUUID)
		isUpdated = true
	}
	if pathced.Date != nil {
		query = query.Set("date = ?", pathced.Date)
		isUpdated = true
	}

	if !isUpdated {
		return nil
	}

	_, err := query.Exec(ctx)
	return err
}

func (r *bunLessonRepository) DeleteLesson(ctx context.Context, lessonUUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.Lesson{UUID: lessonUUID}).WherePK().Exec(ctx)
	return err
}
