package server

import (
	"context"

	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TimetableServer) ListTimetable(ctx context.Context, req *ttpb.ListTimetablesRequest) (*ttpb.ListTimetableResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")

	userGroups, err := s.Client.UsersGroups(ctx, req.UserUuid)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't retrieve user's groups")
	}

	DayToLessons := make(map[ttpb.Day][]*ttpb.ShortConcreteLesson)
	var allLessons []*ttpb.ConcreteLessonShortEntry
	for _, group := range userGroups {
		groupLessons, err := s.timetableRepository.ListTimetable(ctx, group, int(req.WeekNumber), 2026)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't get group lessons for that week")
		}
		for _, lesson := range groupLessons {
			DayToLesson, exists := DayToLessons[ttpb.Day(ttpb.Day_value[string(lesson.Period.DayOfWeek)])]
			if exists {
				DayToLessons[ttpb.Day(ttpb.Day_value[string(lesson.Period.DayOfWeek)])] = append(DayToLesson, &ttpb.ShortConcreteLesson{
					WeekNumber: uint32(lesson.Period.WeekNumber),
					Year:       uint32(lesson.Period.Year),
					Period:     uint32(lesson.Period.Position),
					AbstractLesson: &ttpb.AbstractLesson{
						Uuid: lesson.AbstractUUID,
						Name: lesson.AbstractLesson.Name,
					},
					Building:   lesson.Building,
					Auditorium: lesson.Auditorium,
				})
			} else {
				DayToLessons[ttpb.Day(ttpb.Day_value[string(lesson.Period.DayOfWeek)])] = []*ttpb.ShortConcreteLesson{}
				DayToLessons[ttpb.Day(ttpb.Day_value[string(lesson.Period.DayOfWeek)])] = append(DayToLesson, &ttpb.ShortConcreteLesson{
					WeekNumber: uint32(lesson.Period.WeekNumber),
					Year:       uint32(lesson.Period.Year),
					Period:     uint32(lesson.Period.Position),
					AbstractLesson: &ttpb.AbstractLesson{
						Uuid: lesson.AbstractUUID,
						Name: lesson.AbstractLesson.Name,
					},
					Building:   lesson.Building,
					Auditorium: lesson.Auditorium,
				})
			}

		}
	} //something wicked this way comes

	return &ttpb.ListTimetableResponse{
		Timetable: allLessons,
	}, nil
}

func (s *TimetableServer) GetTimetable(ctx context.Context, req *ttpb.GetTimetableRequest) (*ttpb.GetTimetableResponse, error) {

	userGroups, err := s.Client.UsersGroups(ctx, req.UserUuid)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn't retrieve user's groups")
	}

	allLessons := &ttpb.ConcreteLessonShortEntry{Day: req.Day}
	allLessons.Day = req.Day
	for _, group := range userGroups {
		groupLessons, err := s.timetableRepository.ListTimetable(ctx, group, int(req.WeekNumber), 2026)
		if err != nil {
			return nil, status.Error(codes.Internal, "couldn't get group lessons for that week")
		}
		for _, lesson := range groupLessons {
			allLessons.Lessons = append(allLessons.Lessons, &ttpb.ShortConcreteLesson{
				WeekNumber: uint32(lesson.Period.WeekNumber),
				Year:       uint32(lesson.Period.Year),
				Period:     uint32(lesson.Period.Position),
				AbstractLesson: &ttpb.AbstractLesson{
					Uuid: lesson.AbstractUUID,
					Name: lesson.AbstractLesson.Name,
				},
				Building:   lesson.Building,
				Auditorium: lesson.Auditorium,
			})
		}
	} //something wicked this way comes

	return &ttpb.GetTimetableResponse{
		Timetable: allLessons,
	}, nil
}
