package period

import (
	"context"
	"log/slog"
	"time"

	"github.com/egot3/supel/backend/timetable/internal/logctx"
	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunPeriodRepository struct {
	db *bun.DB
}

func NewPeriodRepository(i do.Injector) (PeriodRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunPeriodRepository{db: db}, nil
}

func (r *bunPeriodRepository) CreatePeriod(ctx context.Context, name string, position int16, startTime, endTime time.Time) error {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "repository"),
	)
	ctx = logctx.WithLogger(ctx, logger)

	_, err := r.db.NewInsert().Model(&models.Period{Name: name, Position: position, StartTime: startTime, EndTime: endTime}).
		Exec(ctx)

	// logger.ErrorContext(ctx, err.Error())
	return err
}

func (r *bunPeriodRepository) Period(ctx context.Context, periodUUID uuid.UUID) (*models.Period, error) {
	period := models.Period{UUID: periodUUID}
	err := r.db.NewSelect().Model(&period).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &period, nil
}

func (r *bunPeriodRepository) PatchPeriod(ctx context.Context, patch models.PeriodPatched) error {
	isUpdated := false

	query := r.db.NewUpdate().Model((*models.Period)(nil)).
		Where("uuid = ?", patch.UUID)
	if patch.Name != nil {
		isUpdated = true
		query = query.Set("name = ?", patch.Name)
	}
	if patch.EndTime != nil {
		isUpdated = true
		query = query.Set("end_time = ?", patch.EndTime)
	}
	if patch.StartTime != nil {
		isUpdated = true
		query = query.Set("start_time = ?", patch.StartTime)
	}
	if patch.Position != nil {
		isUpdated = true
		query = query.Set("position = ?", patch.Position)
	}

	if !isUpdated {
		return nil
	}

	_, err := query.Exec(ctx)
	return err
}

func (r *bunPeriodRepository) DeletePeriod(ctx context.Context, periodUUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.Period{UUID: periodUUID}).WherePK().Exec(ctx)
	return err
}

func (r *bunPeriodRepository) ListPeriods(ctx context.Context, page, size uint32) ([]models.Period, int, error) {
	var periods []models.Period
	total, err := r.db.NewSelect().Model(&periods).Limit(int(size)).Offset(int(page * size)).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return periods, total, nil
}
