package repositories

import (
	"context"

	"github.com/Egot3/supel/backend/timetable/internal/models"
)

type ConcreteLessonRepository interface {
	CreateConcreteLesson(ctx context.Context, cl models.ConcreteLesson) error
	DeleteConcreteLesson(ctx context.Context, uuid string) error
	PatchConcreteLesson(ctx context.Context, patchCl models.PatchConcreteLesson) error
	GetConcreteLesson(ctx context.Context, uuid string) (*models.ConcreteLesson, error)
}
