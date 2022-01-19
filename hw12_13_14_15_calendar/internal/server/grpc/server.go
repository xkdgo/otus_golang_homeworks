//go:generate mkdir -p ./proto
//go:generate protoc EventService.proto --go_out=./proto --go-grpc_out=./proto -I ../../../api

package internalgrpc

import (
	"context"
	"net"
	"time"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/pkg/errors"
	pb "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/server/grpc/proto"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
)

type Logger interface { // TODO
	Log(msg ...interface{})
	Info(msg ...interface{})
	Infof(format string, msg ...interface{})
	Error(msg ...interface{})
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

func NewEventServiceServer(addr string, logger Logger, app Application) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(errors.Wrap(err, ":failed to listen host/port"))
		return nil, err
	}
	grpcServer := grpc.NewServer()
	service := &Service{
		logger: logger,
		app:    app,
	}
	pb.RegisterEventServiceServer(grpcServer, service)
	return &Server{
		router:   grpcServer,
		listener: lis,
		logger:   logger}, nil
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
	id, err := s.app.CreateEvent(ctx, ev.Id, ev.Title, ev.UserID, ev.Datetimestart.AsTime(), ev.Duration.AsDuration(), ev.Alarmtime.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateEventResponse{Id: id}, nil
}

func (s *Service) UpdateEvent(context.Context, *pb.Event) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEvent not implemented")
}

func (s *Service) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (res *emptypb.Empty, err error) {
	err = s.app.DeleteEvent(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Service) ListEventsDay(context.Context, *pb.ListEventsRequest) (*pb.Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsDay not implemented")
}

func (s *Service) ListEventsWeek(context.Context, *pb.ListEventsRequest) (*pb.Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsWeek not implemented")
}

func (s *Service) ListEventsMonth(context.Context, *pb.ListEventsRequest) (*pb.Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEventsMonth not implemented")
}
