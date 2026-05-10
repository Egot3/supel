package repositories

import (
	"context"

	"github.com/Egot3/supel/backend/timetable/internal/models"
	"github.com/Egot3/supel/backend/timetable/internal/types"
)

type PeriodRepository interface {
	Period(ctx context.Context, periodNumber, weekNumber, year uint16, day types.Day) (*models.Period, error)
	PatchPeriod(ctx context.Context, pcpr models.PatchedPeriod) error
}
