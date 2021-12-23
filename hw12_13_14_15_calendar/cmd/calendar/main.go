package main

import (
	"context"
	"fmt"

	// "flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVarP(&configFile, "config", "c", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}
	exitCh := make(chan struct{})
	config := NewConfig()
	logg := logger.New(config.Logger.Level)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		fmt.Println("goroutine started")
		<-ctx.Done()
		fmt.Println("context Done")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		close(exitCh)
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	<-exitCh

}
