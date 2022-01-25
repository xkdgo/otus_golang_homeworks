package app

import (
	"context"
	"sync"
	"time"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	wg      *sync.WaitGroup
	Logger  Logger
	storage storage.Storage
	timeout time.Duration
}

func NewAppScheduler(logger Logger, stor storage.Storage, timeout time.Duration) *Scheduler {
	return &Scheduler{
		wg:      &sync.WaitGroup{},
		Logger:  logger,
		storage: stor,
		timeout: timeout,
	}
}

func (a *Scheduler) Start(ctx context.Context) {
	a.wg.Add(1)
	go a.queryDataToSend(ctx)
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
				a.Logger.Debugf("query data to send error: %q", err)
			}
			for _, event := range events {
				a.Notify(event)
			}
		}
	}
}

func (a *Scheduler) Notify(event storage.Event) {
	a.Logger.Infof("Sended event %v", event)
}
