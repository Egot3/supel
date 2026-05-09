package server

import (
	"context"
	"database/sql"
	"errors"
	"log"

	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"github.com/Egot3/supel/backend/timetable/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TimetableServer) HomeworkAttachmentGETUrl(ctx context.Context, req *ttpb.HomeworkAttachmentGETUrlsRequest) (*ttpb.HomeworkAttachmentGETUrlsResponse, error) {
	log.Println("got request to get get url for attachment")

	keys, err := s.homeworkAttachmentRepository.HomeworkAttachmentKeysByConcreteUUID(ctx, req.GetConcreteLessonUuid())
	if err != nil {
		log.Printf("couldn't get key for %v: %v", req.GetConcreteLessonUuid(), err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Couldn't find a homework attachment")
		}
		return nil, status.Error(codes.Internal, "Couldn't get GET url for homework attachment")
	}

	var urls []string
	for _, key := range keys {
		url, err := s.storageService.GETurl(ctx, key)
		if err != nil {
			log.Printf("couldn't gen a get url for %v: %v ", req.GetConcreteLessonUuid(), err.Error())
			return nil, status.Error(codes.Internal, "Couldn't gen GET url for homework attachment")
		}

		urls = append(urls, url)
	}

	return &ttpb.HomeworkAttachmentGETUrlsResponse{
		HomeworkAttachmentUrls: urls,
	}, nil
}

func (s *TimetableServer) HomeworkAttachmentPUTUrl(ctx context.Context, req *ttpb.HomeworkAttachmentPUTUrlRequest) (*ttpb.HomeworkAttachmentPUTUrlResponse, error) {
	log.Println("got request for put urls for hwa")

	hwa, err := s.homeworkAttachmentRepository.CreateHomeworkAttachment(ctx, models.HomeworkAttachment{
		Name:         req.Name,
		Mime:         req.Mime,
		ConcreteUUID: req.ConcreteLessonUuid,
	})
	if err != nil {
		log.Printf("couldn't crete a key for homework attachment: %v", err.Error())
		return nil, status.Error(codes.Internal, "Couldn't create an attachment")
	}

	url, err := s.storageService.PUTurl(ctx, hwa.StorageKey, hwa.Mime)
	if err != nil {
		log.Printf("couldn't generate a put url: %v", err.Error())
		return nil, status.Error(codes.Internal, "Couldn't gen an attachment")
	}

	return &ttpb.HomeworkAttachmentPUTUrlResponse{
		HomeworkAttachmentUrl: url,
	}, nil
}
