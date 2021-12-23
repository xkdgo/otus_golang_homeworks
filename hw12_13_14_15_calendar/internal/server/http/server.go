package internalhttp

import (
	"context"
	"fmt"
)

type Server struct { // TODO
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	fmt.Println("started StopFunc")
	<-ctx.Done()
	fmt.Println("ctx done ended StopFunc")
	// TODO
	return nil
}

// TODO
