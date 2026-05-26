package period

import (
	"context"
	"time"

	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
)

type PeriodRepository interface {
	CreatePeriod(ctx context.Context, name string, position int16, startTime, endTime time.Time) error
	Period(ctx context.Context, periodUUID uuid.UUID) (*models.Period, error)
	PatchPeriod(ctx context.Context, patch models.PeriodPatched) error
	DeletePeriod(ctx context.Context, periodUUID uuid.UUID) error
}
