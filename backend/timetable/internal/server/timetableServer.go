package server

import (
	"context"

	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TimetableServer) ListTimetable(ctx context.Context, req *ttpb.ListTimetablesRequest) (*ttpb.ListTimetableResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}
