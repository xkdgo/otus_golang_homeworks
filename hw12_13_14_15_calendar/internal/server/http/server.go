package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
}

type Application interface { // TODO
}

func NewServer(addr string, logger Logger, app Application) *Server {
	s := &Server{}
	s.logger = logger
	s.app = app
	s.addr = addr
	return s
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World")
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/", s.loggingMiddleware(http.HandlerFunc(HelloServer)))
	s.router = &http.Server{
		Addr:    s.addr,
		Handler: mux,
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
