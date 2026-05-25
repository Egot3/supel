package repositories

import (
	"context"

	"github.com/Egot3/supel/backend/timetable/internal/models"
	"github.com/Egot3/supel/backend/timetable/internal/types"
)

type TimetableRepository interface {
	ListTimetable(ctx context.Context, groupUUID string, weekNumber, year int) ([]*models.ConcreteLesson, error)
	Timetable(ctx context.Context, groupUUID string, day types.Day, weekNumber, year int) ([]*models.ConcreteLesson, error)
}
