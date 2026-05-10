package repositories

import (
	"context"
	"database/sql"

	"github.com/Egot3/supel/backend/timetable/internal/models"
	"github.com/Egot3/supel/backend/timetable/internal/types"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunPeriodRepository struct {
	db *bun.DB
}

func NewPeriodRepository(i do.Injector) (PeriodRepository, error) {
	db, err := do.Invoke[*bun.DB](i)
	if err != nil {
		return nil, err
	}
	return &bunPeriodRepository{db: db}, nil
}

func (r *bunPeriodRepository) Period(ctx context.Context, periodNumber, weekNumber, year uint16, day types.Day) (*models.Period, error) {
	var pr models.Period

	res, err := r.db.NewSelect().
		Model(&pr).
		Where("position = ?", periodNumber).
		Where("week_number = ?", weekNumber).
		Where("day_of_week = ?", day).
		Where("year = ?", year).
		Exec(ctx)
	if val, _ := res.RowsAffected(); val == 0 {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	return &pr, err
}

func (r *bunPeriodRepository) PatchPeriod(ctx context.Context, pcpr models.PatchedPeriod) error {
	query := r.db.NewUpdate().Model(&models.PatchedPeriod{UUID: pcpr.UUID})
	if pcpr.DayOfWeek != nil {
		query = query.Set("day_of_week = ?", &pcpr.DayOfWeek)
	}
	if pcpr.WeekNumber != nil {
		query = query.Set("week_number = ?", &pcpr.WeekNumber)
	}
	if pcpr.Year != nil {
		query = query.Set("year = ?", &pcpr.Year)
	}
	if pcpr.Start != nil {
		query = query.Set("start = ?", &pcpr.Start)
	}
	if pcpr.End != nil {
		query = query.Set("end = ?", &pcpr.End)
	}

	res, err := query.Exec(ctx)
	if val, _ := res.RowsAffected(); val == 0 {
		return sql.ErrNoRows
	}

	if err != nil {
		return err
	}

	return nil
}
