package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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

const hoursInDay = 24

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
		os.Exit(1)
	}
	period = period.Truncate(time.Second)
	if period < 1*time.Second {
		logg.Error("period should be equal or more than 1s", errors.Wrapf(err, "%s", config.Scheduler.TimeoutQuery))
		cancel()
		os.Exit(1)
	}
	ttlnum, err := strconv.Atoi(config.Scheduler.TTL)
	if err != nil {
		logg.Error("cant parse scheduler ttl:", errors.Wrapf(err, "%s", config.Scheduler.TTL))
		cancel()
		os.Exit(1)
	}
	ttl := time.Duration(ttlnum*hoursInDay) * time.Hour
	reconnectTimeOutInt, err := strconv.Atoi(config.Scheduler.ReconnectTimeOut)
	if err != nil {
		logg.Error("cant parse scheduler reconnectTimeOut:", errors.Wrapf(err, "%s", config.Scheduler.ReconnectTimeOut))
		cancel()
		os.Exit(1)
	}
	reconnectTimeOut := time.Duration(reconnectTimeOutInt) * time.Millisecond

	sender, err := rabbit.NewSender(serviceName, "direct", true, config.Broker.DialString, reconnectTimeOut, logg)
	if err != nil {
		logg.Error("cant init sender:", errors.Wrapf(err, "%s %s", serviceName, config.Broker.DialString))
		cancel()
		os.Exit(1)
	}
	defer sender.Stop()
	scheduler := app.NewAppScheduler(logg, storage, period, ttl, sender, routingKey)
	scheduler.Start(ctx)
}
