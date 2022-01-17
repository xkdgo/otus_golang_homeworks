package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Server struct { // TODO
	logger Logger
	app    Application
	router *http.Server
	addr   string
}

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

	s.logger.Infof("server started on port %s", s.router.Addr)
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
