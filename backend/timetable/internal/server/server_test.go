package server_test

import (
	"context"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	homeworkattachment "github.com/egot3/supel/backend/timetable/internal/database/repositories/homeworkAttachment"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/lesson"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/period"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/subject"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/timetable"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/timetableentry"
	"github.com/egot3/supel/backend/timetable/internal/interceptors"
	"github.com/egot3/supel/backend/timetable/internal/logctx"
	"github.com/egot3/supel/backend/timetable/internal/server"
	testutils "github.com/egot3/supel/backend/timetable/internal/testUtils"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func startProdServer(t *testing.T, i do.Injector) *grpc.ClientConn {
	t.Helper()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer(grpc.UnaryInterceptor(interceptors.LoggingUnaryInterceptor(logger)))

	logctx.WithLogger(t.Context(), logger)

	server, err := server.NewTimetableService(i)
	require.NoError(t, err)

	ttpb.RegisterTimetableServiceServer(srv, server)
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Logf("bufconn server stopped: %v", err)
		}
	}()

	t.Cleanup(func() { srv.GracefulStop(); lis.Close() })

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return conn
}

func TestTimetable_period(t *testing.T) {
	t.Helper()
	t.Parallel()

	i := testutils.NewTestInjector(t)

	do.Provide(i, period.NewPeriodRepository)
	do.Provide(i, lesson.NewLessonRepository)
	do.Provide(i, subject.NewSubjectRepository)
	do.Provide(i, homeworkattachment.NewHomeworkAttachmentRepository)
	do.Provide(i, timetable.NewTimetableRepository)
	do.Provide(i, timetableentry.NewTimetableEntryRepository)

	do.Provide(i, testutils.AllowAllStub)
	do.Provide(i, testutils.AllowAllUsers)

	conn := startProdServer(t, i)
	client := ttpb.NewTimetableServiceClient(conn)

	testCases := []struct {
		desc        string
		errExpected bool
		name        string
		startTime   *time.Time
		endTime     *time.Time
	}{
		{
			desc: "valid all",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
