package server

import (
	"context"

	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"github.com/Egot3/supel/backend/timetable/internal/database/repositories"
	storage "github.com/Egot3/supel/backend/timetable/internal/s3"
	"github.com/samber/do/v2"
	"google.golang.org/grpc/metadata"
)

type TimetableServer struct {
	ttpb.UnimplementedTimetableServiceServer
	storageService               storage.StorageService
	abstractLessonRepository     repositories.AbstractLessonRepository
	concreteLessonRepository     repositories.ConcreteLessonRepository
	homeworkAttachmentRepository repositories.HomeworkAttachmentRepository
	periodRepository             repositories.PeriodRepository
}

func UserFromContext(ctx context.Context) (userID string, role string, ok bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", false
	}
	userID = md.Get("user-uuid")[0]
	role = md.Get("user-role")[0]
	return userID, role, !(len(userID) == 0 && len(role) == 0)
}

func NewTimetableService(i do.Injector) (*TimetableServer, error) {
	storageService, err := do.Invoke[storage.StorageService](i)
	if err != nil {
		return nil, err
	}
	abstractLessonRepository, err := do.Invoke[repositories.AbstractLessonRepository](i)
	if err != nil {
		return nil, err
	}
	concreteLessonRepository, err := do.Invoke[repositories.ConcreteLessonRepository](i)
	if err != nil {
		return nil, err
	}
	homeworkAttachmentRepository, err := do.Invoke[repositories.HomeworkAttachmentRepository](i)
	if err != nil {
		return nil, err
	}
	periodRepository, err := do.Invoke[repositories.PeriodRepository](i)
	if err != nil {
		return nil, err
	}

	return &TimetableServer{
		storageService:               storageService,
		abstractLessonRepository:     abstractLessonRepository,
		concreteLessonRepository:     concreteLessonRepository,
		homeworkAttachmentRepository: homeworkAttachmentRepository,
		periodRepository:             periodRepository,
	}, nil
}
