package timetable

import (
	"context"
	"time"

	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
)

type TimetableRepository interface {
	Timetable(ctx context.Context, timetableUUID uuid.UUID) (*models.Timetable, error)
	CreateTimetable(ctx context.Context, groupUUID uuid.UUID, name string, assignAt, revokeAt *time.Time) error
	PatchTimetable(ctx context.Context, patched models.TimetablePatched) error
	DeleteTimetable(ctx context.Context, timetableUUID uuid.UUID) error
	ListTimetables(ctx context.Context, page, size uint32) ([]models.Timetable, int, error)
	CurrentTimetable(ctx context.Context, groupUUID uuid.UUID) (*models.Timetable, error)
	TimetableByDate(ctx context.Context, groupUUID uuid.UUID, date time.Time) (*models.Timetable, error)
}
