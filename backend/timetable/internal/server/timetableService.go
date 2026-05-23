package server

import (
	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"github.com/Egot3/supel/backend/timetable/internal/database/repositories"
	grpcutils "github.com/Egot3/supel/backend/timetable/internal/grpcUtils"
	storage "github.com/Egot3/supel/backend/timetable/internal/s3"
	"github.com/Egot3/supel/backend/timetable/internal/services"
	"github.com/samber/do/v2"
)

type TimetableServer struct {
	ttpb.UnimplementedTimetableServiceServer
	storageService               storage.StorageService
	abstractLessonRepository     repositories.AbstractLessonRepository
	concreteLessonRepository     repositories.ConcreteLessonRepository
	homeworkAttachmentRepository repositories.HomeworkAttachmentRepository
	periodRepository             repositories.PeriodRepository
	timetableRepository          repositories.TimetableRepository
	Client                       services.Client
	su                           grpcutils.ServerUtils
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

	timetableRepository := do.MustInvoke[repositories.TimetableRepository](i)
	client := do.MustInvoke[services.Client](i)
	su := do.MustInvoke[grpcutils.ServerUtils](i)

	return &TimetableServer{
		storageService:               storageService,
		abstractLessonRepository:     abstractLessonRepository,
		concreteLessonRepository:     concreteLessonRepository,
		homeworkAttachmentRepository: homeworkAttachmentRepository,
		periodRepository:             periodRepository,
		timetableRepository:          timetableRepository,
		Client:                       client,
		su:                           su,
	}, nil
}
