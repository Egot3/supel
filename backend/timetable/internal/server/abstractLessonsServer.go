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
)

// Abstracts
func (s *TimetableServer) CreateAbstractLesson(ctx context.Context, req *ttpb.CreateAbstractLessonRequest) (*emptypb.Empty, error) {
	log.Printf("got request to create Abstract Lesson")
	//RBAC is unimplemented

	err := s.abstractLessonRepository.CreateAbstractLesson(ctx, req.GetName())
	if err != nil {
		log.Printf("Couldn't create abstract lesson: %v", err.Error())
		return nil, status.Error(codes.Internal, "Couldn't create abstract lesson")
	}

	return nil, nil
}

func (s *TimetableServer) DeleteAbstractLesson(ctx context.Context, req *ttpb.DeleteAbstractLessonRequest) (*emptypb.Empty, error) {
	log.Printf("got request to delete Abstract Lesson")
	//RBAC is unimplemented

	err := s.abstractLessonRepository.DeleteAbstractLesson(ctx, req.GetLessonUuid())
	if err != nil {
		log.Printf("Couldn't delete abstract lesson: %v", err.Error())
		return nil, status.Error(codes.Internal, "Couldn't delete abstract lesson")
	}

	return nil, nil
}

func (s *TimetableServer) GetAbstractLesson(ctx context.Context, req *ttpb.GetAbstractLessonRequest) (*ttpb.GetAbstractLessonResponse, error) {
	log.Printf("got request to get Abstract Lesson")

	al, err := s.abstractLessonRepository.GetAbstractLesson(ctx, req.GetLessonUuid())
	if err != nil {
		log.Printf("Abstract lesson getting error: %v", err.Error())
		return nil, status.Error(codes.Internal, "Couldn't fetch abstract lesson")
	}
	if al == nil {
		log.Println("Abstract lesson not found in get")
		return nil, status.Error(codes.NotFound, "Abstract lesson not found")
	}

	return &ttpb.GetAbstractLessonResponse{
		Lesson: &ttpb.AbstractLesson{
			Uuid: al.UUID,
			Name: al.Name,
		},
	}, nil
}

func (s *TimetableServer) PatchAbstractLesson(ctx context.Context, req *ttpb.PatchAbstractLessonRequest) (*emptypb.Empty, error) {
	log.Printf("got request to patch Abstract Lesson")

	err := s.abstractLessonRepository.PatchAbstractLesson(ctx, models.PatchAbstractLesson{UUID: req.GetLessonUuid(), Name: req.Name})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Abstract lesson not found in patch")
			return nil, status.Error(codes.NotFound, "Abstract Lesson not found")
		}
		log.Printf("couldn't fetch Abstract Lesson: %v", err.Error())
		return nil, status.Error(codes.Internal, "Couldn't patch abstract lesson")
	}

	return nil, nil
}
