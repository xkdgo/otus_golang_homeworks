package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	logger Logger
	app    Application
	router *http.Server
	addr   string
}

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

func NewServer(addr string, logger Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
		addr:   addr,
	}
}

func (s *Server) Start(ctx context.Context) error {
	handler := NewRootHandler(s.app, s.logger)
	s.router = &http.Server{
		Addr:    s.addr,
		Handler: s.loggingMiddleware(s.authMiddleware(handler)),
	}

	s.logger.Infof("http server started on port %s", s.router.Addr)
	err := s.router.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	fmt.Println("started StopFunc")
	if err := s.router.Shutdown(ctx); err != nil {
		fmt.Printf("%#v\n", err)
		return err
	}
	<-ctx.Done()
	fmt.Println("ctx done ended StopFunc")
	// TODO
	return nil
}
