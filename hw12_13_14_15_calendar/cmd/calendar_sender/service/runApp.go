package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/config"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue/controllers/rabbit"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/plugins/logger/zap"
)

var exchangeName = serviceNameExchange

func RunApp(config config.SenderConfig) {
	pluginlogger, err := zap.NewLogger(logger.WithFields(map[string]interface{}{serviceName: ""}))
	if err != nil {
		fmt.Println("Cant initialize zap logger")
		os.Exit(1)
	}
	logg := logger.New(config.Logger.Level, pluginlogger)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	reconnectTimeOutInt, err := strconv.Atoi(config.Sender.ReconnectTimeOut)
	if err != nil {
		logg.Fatalf("cant parse scheduler reconnectTimeOut:", errors.Wrapf(err, "%s", config.Sender.ReconnectTimeOut))
	}
	reconnectTimeOut := time.Duration(reconnectTimeOutInt) * time.Millisecond

	sender := app.NewAppSender(logg,
		config.Sender.NumWorkers,
		config.Broker.DialString,
		exchangeName,
		config.Sender.RoutingKey,
		serviceName,
		reconnectTimeOut,
		&rabbit.Receiver{})
	err = sender.Init()
	if err != nil {
		logg.Fatalf("error during init application %v", err)
	}
	sender.Start(ctx)
}
