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
	wg         *sync.WaitGroup
	Logger     Logger
	storage    storage.Storage
	timeout    time.Duration
	ttl        time.Duration
	Notifier   queue.Notifier
	routingKey string
}

func NewAppScheduler(
	logger Logger,
	stor storage.Storage,
	timeout, ttl time.Duration,
	notifier queue.Notifier,
	routingKey string) *Scheduler {
	return &Scheduler{
		wg:         &sync.WaitGroup{},
		Logger:     logger,
		storage:    stor,
		timeout:    timeout,
		ttl:        ttl,
		Notifier:   notifier,
		routingKey: routingKey,
	}
}

func (a *Scheduler) Start(ctx context.Context) {
	a.wg.Add(3)
	go a.queryDataToSend(ctx)
	go a.monitorConnectToQueue()
	go a.cleanOldEvents(ctx)
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
			timeStart := time.Now().UTC().Truncate(time.Second)
			periodTimeEnd := timeStart.Add(a.timeout).Truncate(time.Second)
			a.Logger.Debugf("querydatatosend period: %s to %s", timeStart, periodTimeEnd)
			events, err := a.storage.ListEventsToNotify(timeStart, periodTimeEnd)
			if err != nil {
				a.Logger.Debugf("querydatatosend error: %q", err)
			}
			for _, event := range events {
				err := a.notify(event)
				if err != nil {
					a.Logger.Errorf("cant send event to queue %q", err)
				}
			}
		}
	}
}

func (a *Scheduler) notify(event storage.Event) error {
	var m queue.NotifyEvent
	m.ID = event.ID
	m.UserID = event.UserID
	m.DateTimeStart = event.DateTimeStart.UTC().Format(time.RFC3339)
	m.Title = event.Title
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = a.Notifier.Publish(a.routingKey, "application/json", body)
	if err != nil {
		return err
	}
	a.Logger.Debugf("sended event %v", string(body))
	return nil
}

func (a *Scheduler) cleanOldEvents(ctx context.Context) {
	defer a.wg.Done()
	timer := time.NewTicker(10 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			events, err := a.storage.ListEventsToDelete(a.ttl)
			if err != nil {
				a.Logger.Debugf("querydatatodeleteerror: %q", err)
			}
			for _, event := range events {
				err := a.storage.DeleteEvent(event.ID)
				if err != nil {
					a.Logger.Errorf("querydatatodeleteerror: %q", err)
				} else {
					a.Logger.Debugf("deleted event id:%v title: '%v'", event.ID, event.Title)
				}
			}
		}
	}
}

func (a *Scheduler) monitorConnectToQueue() {
	defer a.wg.Done()
	go a.Notifier.Listen()
}
