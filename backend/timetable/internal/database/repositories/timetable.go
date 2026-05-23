package repositories

import (
	"context"

	"github.com/Egot3/supel/backend/timetable/internal/models"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunTimetableRepository struct {
	db *bun.DB
}

func NewTimetableRepository(i do.Injector) (TimetableRepository, error) {
	db, err := do.Invoke[*bun.DB](i)
	if err != nil {
		return nil, err
	}
	return &bunTimetableRepository{db: db}, nil
}

func (r *bunTimetableRepository) ListTimetable(ctx context.Context, groupUUID string, weekNumber, year int) ([]*models.ConcreteLesson, error) {
	var lessons []*models.ConcreteLesson
	err := r.db.NewSelect().Model(&lessons).
		Where("group_uuid = ?", groupUUID).Where("year = ?", year).Where("week_number = ?", weekNumber).
		Order("day_of_week ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return lessons, nil
}
