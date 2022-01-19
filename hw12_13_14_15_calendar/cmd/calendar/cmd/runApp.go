package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/helper"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/server/http"
)

func RunApp(config Config) {
	exitCh := make(chan struct{})
	logg := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := helper.InitStorage(config.Storage.Type, config.Storage.DSN)
	if err != nil {
		logg.Error("cant init storage:", errors.Wrapf(err, "%s", config.Storage.Type))
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	defer storage.Close()
	calendar := app.New(logg, storage)

	serverHTTP := internalhttp.NewServer(
		net.JoinHostPort(config.ServerHTTP.Host, config.ServerHTTP.Port),
		logg,
		calendar)
	listener, err := internalgrpc.GetListener(
		net.JoinHostPort(config.ServerGRPC.Host, config.ServerGRPC.Port),
	)
	if err != nil {
		logg.Error(errors.Wrap(err, ":failed to listen host/port"))
		cancel()
		os.Exit(1)
	}
	serverGRPC, err := internalgrpc.NewEventServiceServer(
		listener,
		logg,
		calendar)
	if err != nil {
		logg.Error(
			"cant init serviceGRPC:",
			errors.Wrapf(err, "%s", net.JoinHostPort(config.ServerHTTP.Host, config.ServerHTTP.Port)))
		cancel()
		os.Exit(1)
	}

	go func() {
		fmt.Println("listen to stop signal goroutine started")
		<-ctx.Done()
		fmt.Println("context Done")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := serverHTTP.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		if err := serverGRPC.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
		close(exitCh)
	}()

	logg.Info("calendar is running...")
	go func() {
		if err := serverHTTP.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()
	go func() {
		if err := serverGRPC.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()

	<-exitCh
}
