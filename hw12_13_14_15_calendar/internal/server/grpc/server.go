//go:generate mkdir -p ./proto
//go:generate protoc EventService.proto --go_out=./proto --go-grpc_out=./proto -I ../../../api

package internalgrpc

import (
	"context"
	"net"
	"time"

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
		Duration, AlarmTime time.Duration) (createdID string, err error)
	UpdateEvent(ctx context.Context,
		id, title, userID string,
		dateTimeStart time.Time,
		duration, alarmTime time.Duration) (err error)
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
