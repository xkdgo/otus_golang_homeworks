package app

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/queue"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	wg       *sync.WaitGroup
	Logger   Logger
	storage  storage.Storage
	timeout  time.Duration
	Notifier queue.Notifier
}

func NewAppScheduler(logger Logger, stor storage.Storage, timeout time.Duration, notifier queue.Notifier) *Scheduler {
	return &Scheduler{
		wg:       &sync.WaitGroup{},
		Logger:   logger,
		storage:  stor,
		timeout:  timeout,
		Notifier: notifier,
	}
}

func (a *Scheduler) Start(ctx context.Context) {
	a.wg.Add(2)
	go a.queryDataToSend(ctx)
	go a.MonitorConnectToQueue()
	a.wg.Wait()
}

func (a *Scheduler) Stop() {
	a.Logger.Debugf("stop calendar scheduler")
}

func (a *Scheduler) queryDataToSend(ctx context.Context) {
	defer a.wg.Done()
	timer := time.NewTicker(a.timeout)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			a.Stop()
			return
		case <-timer.C:
			timeStart := time.Now().Truncate(time.Second)
			periodTimeEnd := timeStart.Add(a.timeout).Truncate(time.Second)
			events, err := a.storage.ListEventsToNotify(timeStart, periodTimeEnd)
			if err != nil {
				a.Logger.Debugf("querydatatosend error: %q", err)
			}
			for _, event := range events {
				err := a.Notify(event)
				if err != nil {
					a.Logger.Errorf("cant send event to queue %q", err)
				}
			}
		}
	}
}

func (a *Scheduler) Notify(event storage.Event) error {
	var m queue.NotifyEvent
	m.ID = event.ID
	m.UserID = event.UserID
	m.DateTimeStart = event.DateTimeStart.Format(time.RFC3339)
	m.Title = event.Title
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = a.Notifier.Publish("calendar_sender", "application/json", body)
	if err != nil {
		return err
	}
	return nil
}

func (a *Scheduler) MonitorConnectToQueue() {
	defer a.wg.Done()
	go a.Notifier.Listen()
}