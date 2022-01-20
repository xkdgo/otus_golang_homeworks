package internalgrpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/helper"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/server/grpc/proto"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/plugins/logger/zap"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	status "google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	pluginlogger, err := zap.NewLogger()
	if err != nil {
		return
	}
	logg := logger.New("DEBUG", pluginlogger)
	lis = bufconn.Listen(bufSize)
	storage, err := helper.InitStorage("in-memory", "")
	if err != nil {
		return
	}
	calendar := app.New(logg, storage)
	s, err := NewEventServiceServer(lis,
		logg,
		calendar)
	if err != nil {
		return
	}
	go func() {
		if err := s.router.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGrpcAPI(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()
	client := pb.NewEventServiceClient(conn)
	id := "0fb2cc2e-85c5-41e2-b83e-66e961cf43db"
	userID := "db46a463-73dc-47cd-ac0e-886b2a99981a"
	in := &pb.Event{
		Id:     id,
		Title:  "Hello",
		UserID: userID,
		Datetimestart: &timestamppb.Timestamp{
			Seconds: 1643717700,
			Nanos:   10,
		},
		Duration: &durationpb.Duration{
			Seconds: 3600,
			Nanos:   10,
		},
		Alarmtime: &timestamppb.Timestamp{
			Seconds: 1643716800,
			Nanos:   10,
		},
	}
	resp, err := client.CreateEvent(ctx, in)
	require.NoError(t, err)
	require.Equal(t, id, resp.Id)

	in.Title = "Updated Title"
	updatedResonse, err := client.UpdateEvent(ctx, in)
	require.NoError(t, err)
	require.IsType(t, &emptypb.Empty{}, updatedResonse)

	listEvRequest := &pb.ListEventsRequest{UserID: userID, Datetimestart: "2022-02-01"}

	listresp, err := client.ListEventsDay(ctx, listEvRequest)
	require.NoError(t, err)
	require.Len(t, listresp.Events, 1)
	require.True(t, proto.Equal(in, listresp.Events[0]))

	listresp, err = client.ListEventsWeek(ctx, listEvRequest)
	require.NoError(t, err)
	require.Len(t, listresp.Events, 1)
	require.True(t, proto.Equal(in, listresp.Events[0]))

	listresp, err = client.ListEventsMonth(ctx, listEvRequest)
	require.NoError(t, err)
	require.Len(t, listresp.Events, 1)
	require.True(t, proto.Equal(in, listresp.Events[0]))

	listEvBadRequest := &pb.ListEventsRequest{UserID: userID, Datetimestart: "Invalid Argument"}
	_, err = client.ListEventsDay(ctx, listEvBadRequest)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())

	deleteResonse, err := client.DeleteEvent(ctx, &pb.DeleteEventRequest{Id: in.Id})
	require.NoError(t, err)
	require.IsType(t, &emptypb.Empty{}, deleteResonse)

	listresp, err = client.ListEventsMonth(ctx, listEvRequest)
	require.NoError(t, err)
	require.Len(t, listresp.Events, 0)
}
