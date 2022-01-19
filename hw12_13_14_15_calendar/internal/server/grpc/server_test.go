package internalgrpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/helper"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/server/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	logg := logger.New("DEBUG")
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
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventServiceClient(conn)
	in := &pb.Event{
		Id:     "0fb2cc2e-85c5-41e2-b83e-66e961cf43db",
		Title:  "Hello",
		UserID: "db46a463-73dc-47cd-ac0e-886b2a99981a",
		Datetimestart: &timestamppb.Timestamp{
			Seconds: 20,
			Nanos:   10,
		},
		Duration: &durationpb.Duration{
			Seconds: 20,
			Nanos:   10,
		},
		Alarmtime: &timestamppb.Timestamp{
			Seconds: 20,
			Nanos:   10,
		},
	}
	resp, err := client.CreateEvent(ctx, in)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}
	log.Printf("Response: %+v", resp)
	// Test for output here.
}
