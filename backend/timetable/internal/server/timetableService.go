package server

import (
	"context"

	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"github.com/Egot3/supel/backend/timetable/internal/database/repositories"
	storage "github.com/Egot3/supel/backend/timetable/internal/s3"
	"google.golang.org/grpc/metadata"
)

type TimetableServer struct {
	ttpb.UnimplementedTimetableServiceServer
	storageService               storage.StorageService
	abstractLessonRepository     repositories.AbstractLessonRepository
	concreteLessonRepository     repositories.ConcreteLessonRepository
	homeworkAttachmentRepository repositories.HomeworkAttachmentRepository
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

func NewTimetableService(storageService storage.StorageService, abstractLessonRepository repositories.AbstractLessonRepository, concreteLessonRepository repositories.ConcreteLessonRepository, homeworkAttachmentRepository repositories.HomeworkAttachmentRepository) *TimetableServer {
	return &TimetableServer{
		storageService:               storageService,
		abstractLessonRepository:     abstractLessonRepository,
		concreteLessonRepository:     concreteLessonRepository,
		homeworkAttachmentRepository: homeworkAttachmentRepository,
	}
}
