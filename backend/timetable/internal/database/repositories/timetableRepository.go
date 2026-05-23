package repositories

import (
	"context"

	"github.com/Egot3/supel/backend/timetable/internal/models"
)

type TimetableRepository interface {
	ListTimetable(ctx context.Context, groupUUID string, weekNumber, year int) ([]*models.ConcreteLesson, error)
}
