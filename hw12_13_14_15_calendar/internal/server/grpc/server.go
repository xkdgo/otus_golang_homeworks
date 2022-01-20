//go:generate mkdir -p ./proto
//go:generate protoc EventService.proto --go_out=./proto --go-grpc_out=./proto -I ../../../api

package internalgrpc

import (
	"context"
	"net"
	"time"

	pb "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/server/grpc/proto"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	timelayout = "2006-01-02"
)

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Trace(args ...interface{})
	Tracef(template string, args ...interface{})
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Application interface {
	CreateEvent(ctx context.Context,
		id, title, userID string,
		DateTimeStart time.Time,
		Duration time.Duration,
		AlarmTime time.Time) (createdID string, err error)
	UpdateEvent(ctx context.Context,
		id, title, userID string,
		dateTimeStart time.Time,
		duration time.Duration,
		alarmTime time.Time) (err error)
	DeleteEvent(ctx context.Context, id string) error
	ListEventsDay(ctx context.Context, userID string, dateTime time.Time) (events []storage.Event, err error)
	ListEventsWeek(ctx context.Context, userID string, dateTime time.Time) (events []storage.Event, err error)
	ListEventsMonth(ctx context.Context, userID string, dateTime time.Time) (events []storage.Event, err error)
	GetStorage() storage.Storage
}

type Service struct {
	pb.UnimplementedEventServiceServer
	logger Logger
	app    Application
}

type Server struct {
	router   *grpc.Server
	listener net.Listener
	logger   Logger
}

func GetListener(addr string) (net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return lis, nil
}

func NewEventServiceServer(lis net.Listener, logger Logger, app Application) (*Server, error) {
	grpcServer := grpc.NewServer()
	service := &Service{
		logger: logger,
		app:    app,
	}
	pb.RegisterEventServiceServer(grpcServer, service)
	return &Server{
		router:   grpcServer,
		listener: lis,
		logger:   logger,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Infof("GRPC server listening on %v", s.listener.Addr())
	err := s.router.Serve(s.listener)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.router.GracefulStop()
	return nil
}

func (s *Service) CreateEvent(ctx context.Context, ev *pb.Event) (res *pb.CreateEventResponse, err error) {
	id, err := s.app.CreateEvent(
		ctx,
		ev.Id,
		ev.Title,
		ev.UserID,
		ev.Datetimestart.AsTime(),
		ev.Duration.AsDuration(),
		ev.Alarmtime.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateEventResponse{Id: id}, nil
}

func (s *Service) UpdateEvent(ctx context.Context, ev *pb.Event) (empty *emptypb.Empty, err error) {
	err = s.app.UpdateEvent(
		ctx,
		ev.Id,
		ev.Title,
		ev.UserID,
		ev.Datetimestart.AsTime(),
		ev.Duration.AsDuration(),
		ev.Alarmtime.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (res *emptypb.Empty, err error) {
	err = s.app.DeleteEvent(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) ListEventsDay(ctx context.Context, req *pb.ListEventsRequest) (res *pb.Events, err error) {
	day := req.Datetimestart
	dayTime, err := time.Parse(timelayout, day)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	events, err := s.app.ListEventsDay(ctx, req.UserID, dayTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	protoEvents := convertToProtoEvents(events)
	return &pb.Events{Events: protoEvents}, nil
}

func (s *Service) ListEventsWeek(ctx context.Context, req *pb.ListEventsRequest) (res *pb.Events, err error) {
	day := req.Datetimestart
	dayTime, err := time.Parse(timelayout, day)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	events, err := s.app.ListEventsWeek(ctx, req.UserID, dayTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	protoEvents := convertToProtoEvents(events)
	return &pb.Events{Events: protoEvents}, nil
}

func (s *Service) ListEventsMonth(ctx context.Context, req *pb.ListEventsRequest) (res *pb.Events, err error) {
	day := req.Datetimestart
	dayTime, err := time.Parse(timelayout, day)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	events, err := s.app.ListEventsWeek(ctx, req.UserID, dayTime)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	protoEvents := convertToProtoEvents(events)
	return &pb.Events{Events: protoEvents}, nil
}

func convertToProtoEvents(events []storage.Event) (protoEvents []*pb.Event) {
	protoEvents = make([]*pb.Event, 0, len(events))
	for _, event := range events {
		protoEvents = append(protoEvents, &pb.Event{
			Id:            event.ID,
			Title:         event.Title,
			UserID:        event.UserID,
			Datetimestart: timestamppb.New(event.DateTimeStart),
			Duration:      durationpb.New(event.Duration),
			Alarmtime:     timestamppb.New(event.AlarmTime),
		})
	}
	return protoEvents
}
