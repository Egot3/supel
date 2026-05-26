package server

import (
	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	homeworkattachment "github.com/egot3/supel/backend/timetable/internal/database/repositories/homeworkAttachment"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/lesson"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/period"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/subject"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/timetable"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/timetableentry"
	grpcutils "github.com/egot3/supel/backend/timetable/internal/grpcUtils"
	storage "github.com/egot3/supel/backend/timetable/internal/s3"
	"github.com/egot3/supel/backend/timetable/internal/services"
	"github.com/samber/do/v2"
)

type TimetableService struct {
	ttpb.UnimplementedTimetableServiceServer

	periodRepository             period.PeriodRepository
	subjectRepository            subject.SubjectRepository
	timetableRepository          timetable.TimetableRepository
	timetableEntryRepository     timetableentry.TimetableEntryRepository
	lessonRepository             lesson.LessonRepository
	homeworkAttachmentRepository homeworkattachment.HomeworkAttachmentRepository

	storage storage.StorageService

	grpcClient services.Client
	su         grpcutils.ServerUtils
}

func NewPuddleService(i do.Injector) (TimetableService, error) {
	peRepo := do.MustInvoke[period.PeriodRepository](i)
	suRepo := do.MustInvoke[subject.SubjectRepository](i)
	ttRepo := do.MustInvoke[timetable.TimetableRepository](i)
	teRepo := do.MustInvoke[timetableentry.TimetableEntryRepository](i)
	leRepo := do.MustInvoke[lesson.LessonRepository](i)
	haRepo := do.MustInvoke[homeworkattachment.HomeworkAttachmentRepository](i) //all are named by their initials

	storage := do.MustInvoke[storage.StorageService](i)

	grpcClient := do.MustInvoke[services.Client](i)
	su := do.MustInvoke[grpcutils.ServerUtils](i)
	return TimetableService{
		periodRepository:             peRepo,
		subjectRepository:            suRepo,
		timetableRepository:          ttRepo,
		timetableEntryRepository:     teRepo,
		lessonRepository:             leRepo,
		homeworkAttachmentRepository: haRepo,
		storage:                      storage,

		grpcClient: grpcClient,
		su:         su,
	}, nil
}
