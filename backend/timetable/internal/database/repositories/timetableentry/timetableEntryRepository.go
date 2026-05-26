package timetableentry

import (
	"context"

	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
)

type TimetableEntryRepository interface {
	TimetableEntry(ctx context.Context, timetableEntryUUID uuid.UUID) (*models.TimetableEntry, error)
	TimetableEntriesByTimetable(ctx context.Context, timetableUUID uuid.UUID) ([]models.TimetableEntry, error)
	CreateTimetableEntry(ctx context.Context, timetableUUID, periodUUID, subjectUUID, teacherUUID uuid.UUID, dayOfWeek int16, place string) error
	PatchTimetableEntry(ctx context.Context, patched models.TimetableEntryPatched) error
	DeleteTimetableEntry(ctx context.Context, timetableEntryUUID uuid.UUID) error
}
