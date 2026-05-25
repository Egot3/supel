package server

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"github.com/Egot3/supel/backend/timetable/internal/models"
	"github.com/Egot3/supel/backend/timetable/internal/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *TimetableServer) GetPeriods(ctx context.Context, req *ttpb.GetPeriodRequest) (*ttpb.GetPeriodsResponse, error) {
	periods, err := s.periodRepository.PeriodsByDay(ctx, uint16(req.WeekNumber), uint16(req.Year), types.ProtoDayToDay(req.Day))
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

func (s *TimetableServer) PatchPeriod(ctx context.Context, req *ttpb.PatchPeriodRequest) (*emptypb.Empty, error) {

	ownUUID, ok := s.su.UserFromContext(ctx)
	if !ok {
		log.Panicf("uuid is not ok: %v", ownUUID)
		return nil, status.Error(codes.InvalidArgument, "bad uuid")
	}

	can, err := s.Client.HasPermission(ctx, ownUUID, "periods", nil, rbacpb.Verb_POST)
	if err != nil {
		return nil, status.Error(codes.Internal, "rbac got an error while getting permissions")
	}
	if !can {
		return nil, status.Error(codes.PermissionDenied, "you lack permissions")
	}

	var start *time.Time = nil
	if st := req.Start; st != nil {
		tmp := st.AsTime()
		start = &tmp
	}

	var end *time.Time = nil
	if en := req.End; en != nil {
		tmp := en.AsTime()
		end = &tmp
	}

	var weekNumber *uint16 = nil
	if wn := req.WeekNumber; wn != nil {
		tmp := uint16(*wn)
		weekNumber = &tmp
	}

	var year *uint16 = nil
	if yr := req.Year; yr != nil {
		tmp := uint16(*yr)
		year = &tmp
	}

	var dayOfWeek *types.Day = nil
	if day := req.Day; day != nil {
		tmp := types.Day(*day)
		dayOfWeek = &tmp
	}

	patchedPeriod := models.PatchedPeriod{
		UUID:       req.PeriodUuid,
		Start:      start,
		End:        end,
		WeekNumber: weekNumber,
		Year:       year,
		DayOfWeek:  dayOfWeek,
	}

	err = s.periodRepository.PatchPeriod(ctx, patchedPeriod)
	if err != nil {
		log.Printf("couldn't patch period: %v", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "period to patch was not found")
		}
		return nil, status.Error(codes.Internal, "couldn't patch period")
	}

	return nil, nil
}
