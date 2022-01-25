package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/config"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/helper"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue/controllers/rabbit"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/plugins/logger/zap"
)

func RunApp(config config.SchedulerConfig) {
	// exitCh := make(chan struct{})
	pluginlogger, err := zap.NewLogger(logger.WithFields(map[string]interface{}{"scheduler": ""}))
	if err != nil {
		fmt.Println("Cant initialize zap logger")
		os.Exit(1)
	}
	logg := logger.New(config.Logger.Level, pluginlogger)

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := helper.InitStorage(config.Storage.Type, config.Storage.DSN)
	if err != nil {
		logg.Error("cant init storage:", errors.Wrapf(err, "%s", config.Storage.Type))
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	period, err := time.ParseDuration(config.Scheduler.TimeoutQuery)
	if err != nil {
		logg.Error("cant parse scheduler period:", errors.Wrapf(err, "%s", config.Scheduler.TimeoutQuery))
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	if period < time.Duration(1*time.Second) {
		logg.Error("period should be equal or more than 1s", errors.Wrapf(err, "%s", config.Scheduler.TimeoutQuery))
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	//TODO make amqp config
	sender, err := rabbit.NewSender(serviceName, "direct", true, "amqp://guest:guest@localhost:5672/", logg)
	if err != nil {
		logg.Error("cant init sender:", errors.Wrapf(err, "%s %s", serviceName, "amqp://guest:guest@localhost:5672/"))
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	defer sender.Stop()
	scheduler := app.NewAppScheduler(logg, storage, period, sender)
	scheduler.Start(ctx)
}
