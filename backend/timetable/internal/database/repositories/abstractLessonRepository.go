package repositories

import (
	"context"

	"github.com/Egot3/supel/backend/timetable/internal/models"
)

type AbstractLessonRepository interface {
	CreateAbstractLesson(ctx context.Context, name string) error
	DeleteAbstractLesson(ctx context.Context, uuid string) error
	GetAbstractLesson(ctx context.Context, uuid string) (*models.AbstractLesson, error)
	PatchAbstractLesson(ctx context.Context, patchedAl models.PatchAbstractLesson) error
}
