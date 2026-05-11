package server

import (
	"context"
	"log"

	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"github.com/Egot3/supel/backend/timetable/internal/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *TimetableServer) GetPeriods(ctx context.Context, req *ttpb.GetPeriodRequest) (*ttpb.GetPeriodsResponse, error) {
	periods, err := s.periodRepository.PeriodsByDay(ctx, uint16(req.WeekNumber), uint16(req.Year), types.Day(req.Day))
	if err != nil {
		log.Printf("couldn't select all periods for day: %v", err)
		return nil, status.Error(codes.Internal, "couldn't get all periods for the day")
	}

	protoPeriods := make([]*ttpb.Period, periods[len(periods)-1].Position+1)
	for i, period := range periods {
		protoPeriods[i] = &ttpb.Period{
			Number: uint32(period.Position),
			Start:  timestamppb.New(period.Start),
			End:    timestamppb.New(period.End),
		}
	}

	return &ttpb.GetPeriodsResponse{
		Periods: protoPeriods,
	}, nil
}
