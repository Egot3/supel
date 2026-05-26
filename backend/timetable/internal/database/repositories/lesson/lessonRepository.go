package lesson

import (
	"context"
	"time"

	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
)

type LessonRepository interface {
	Lesson(ctx context.Context, lessonUUID uuid.UUID) (*models.Lesson, error)
	CreateLesson(ctx context.Context, timetableEntryUUID uuid.UUID, date time.Time) error
	PatchLesson(ctx context.Context, pathced models.LessonPatched) error
	DeleteLesson(ctx context.Context, lessonUUID uuid.UUID) error
	LessonByEntry(ctx context.Context, entryUUId uuid.UUID) (*models.Lesson, error)
}
