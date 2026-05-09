package server

import (
	"context"

	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
)

func (s *TimetableServer) ListTimetable(ctx context.Context, req *ttpb.ListTimetablesRequest) (*ttpb.ListTimetableResponse, error)
