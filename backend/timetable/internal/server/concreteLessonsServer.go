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
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *TimetableServer) CreateConcreteLesson(ctx context.Context, req *ttpb.CreateConcreteLessonRequest) (*emptypb.Empty, error) {
	log.Println("Got request to create concrete Lesson")

	period := uint16(req.Period) //ну кто бы мог подумать
	start := req.Start.AsTime()
	end := req.End.AsTime()

	err := s.concreteLessonRepository.CreateConcreteLesson(ctx, models.ConcreteLesson{
		TeacherUUID:  req.GetTeacherUuid(),
		GroupUUID:    req.GetGroupUuid(),
		AbstractUUID: req.GetAbstractLessonUuid(),

		StartTime: &start,
		EndTime:   &end,
		Period:    &period,
	})
	if err != nil {
		log.Printf("Couldn't create concrete lesson: %v", err.Error())
		return nil, status.Error(codes.Internal, "couldn't create concrete lesson")
	}

	return nil, nil
}

func (s *TimetableServer) DeleteConcreteLesson(ctx context.Context, req *ttpb.DeleteConcreteLessonRequest) (*emptypb.Empty, error) {
	log.Println("got request to delete concrete lesson")

	err := s.concreteLessonRepository.DeleteConcreteLesson(ctx, req.GetConcreteLessonUuid())
	if err != nil {
		log.Printf("couldn't delete concrete lesson: %v", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "concrete lesson to delete was not found")
		}
		return nil, status.Error(codes.Internal, "couldn't delete concrete lesson")
	}

	return nil, nil
}

func (s *TimetableServer) GetConcreteLesson(ctx context.Context, req *ttpb.GetConcreteLessonRequest) (*ttpb.GetConcreteLessonResponse, error) {
	log.Println("Got request to get concrete lesson")

	concreteLesson, err := s.concreteLessonRepository.GetConcreteLesson(ctx, req.GetConcreteLessonUuid())
	if err != nil {
		log.Printf("couldn't get concrete lesson: %v", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "concrete lesson to get was not found")
		}
		return nil, status.Error(codes.Internal, "Coudn't get concrete lesson")
	}

	var bodyURL *string = nil
	b, err := s.storageService.GETurl(ctx, concreteLesson.HomeworkBodyKey)
	if err == nil {
		bodyURL = &b
	}

	var attachmentUrls []string
	for _, keyModel := range concreteLesson.Attachments {
		attachmentUrl, err := s.storageService.GETurl(ctx, keyModel.StorageKey)
		if err != nil {
			log.Printf("couldn't generate attachment url for %v: %v", concreteLesson.ConcreteUUID, err.Error())
			return nil, status.Error(codes.Internal, "Couldn't generate attachment urls")
		}
		attachmentUrls = append(attachmentUrls, attachmentUrl)
	}

	lesson := &ttpb.ConcreteLesson{ //doesn't return marks for now
		LessonInfo: &ttpb.AbstractLesson{
			Uuid: concreteLesson.AbstractLesson.UUID,
			Name: concreteLesson.AbstractLesson.Name,
		},
		TeacherUuid:                concreteLesson.TeacherUUID,
		LessonUuid:                 concreteLesson.ConcreteUUID,
		HomeworkTextGetUrl:         bodyURL,
		HomeworkAttachmentsGetUrls: attachmentUrls,
	}
	if concreteLesson.StartTime != nil {
		lesson.Start = timestamppb.New(*concreteLesson.StartTime)
	}
	if concreteLesson.EndTime != nil {
		lesson.End = timestamppb.New(*concreteLesson.EndTime)
	}
	if concreteLesson.Period != nil {
		lesson.Period = uint32(*concreteLesson.Period)
	}
	return &ttpb.GetConcreteLessonResponse{
		Lesson: lesson,
	}, nil
}

func (s *TimetableServer) PatchConcreteLesson(ctx context.Context, req *ttpb.PatchConcreteLessonRequest) (*emptypb.Empty, error) {
	log.Println("Got request to patch concrete Lesson")
	//RBAC here

	start := req.Start.AsTime()
	end := req.End.AsTime()
	err := s.concreteLessonRepository.PatchConcreteLesson(ctx, models.PatchConcreteLesson{
		ConcreteUUID: req.ConcreteLessonUuid,
		AbstractUUID: req.AbstractLessonUuid,
		TeacherUUID:  req.TeacherUuid,
		GroupUUID:    req.GroupUuid,
		StartTime:    &start,
		EndTime:      &end,
	})
	if err != nil {
		log.Printf("couldn't patch concrete lesson: %v", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "concrete lesson to patch was not found")
		}
		return nil, status.Error(codes.Internal, "Coudn't patch concrete lesson")
	}

	return nil, nil
}
