package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	wg           *sync.WaitGroup
	mu           *sync.Mutex
	Logger       Logger
	storage      storage.Storage
	timeout      time.Duration
	notfyQueue   chan storage.Event
	schedulerMap map[time.Time][]storage.Event
}

func NewAppScheduler(logger Logger, stor storage.Storage, timeout time.Duration) *Scheduler {
	return &Scheduler{
		wg:           &sync.WaitGroup{},
		mu:           &sync.Mutex{},
		Logger:       logger,
		storage:      stor,
		timeout:      timeout,
		notfyQueue:   make(chan storage.Event, 1),
		schedulerMap: make(map[time.Time][]storage.Event)}
}

func (a *Scheduler) Start(ctx context.Context) {
	a.wg.Add(1)
	go a.queryDataToSend(ctx)
	a.wg.Wait()
}

func (a *Scheduler) Stop() {
	a.Logger.Debugf("started stop calendar scheduler")
	close(a.notfyQueue)
}

func (a *Scheduler) fillSchedulerMap() (err error) {
	timeStart := time.Now().Truncate(time.Second)
	periodTimeEnd := timeStart.Add(a.timeout).Truncate(time.Second)
	events, err := a.storage.ListEventsToNotify(timeStart, periodTimeEnd)
	if err != nil {
		return err
	}
	for _, event := range events {
		fmt.Println("key = ", timeStart)
		a.mu.Lock()
		a.schedulerMap[timeStart] = append(a.schedulerMap[timeStart], event)
		a.mu.Unlock()
		a.Notify(event)
	}
	return nil
}

func (a *Scheduler) queryDataToSend(ctx context.Context) {
	defer a.wg.Done()
	timer := time.NewTicker(a.timeout)
	for {
		select {
		case <-ctx.Done():
			close(a.notfyQueue)
			return
		case <-timer.C:
			err := a.fillSchedulerMap()
			if err != nil {
				a.Logger.Debugf("query data to send error: %q", err)

			}
		}
	}
}

func (a *Scheduler) Notify(event storage.Event) {
	a.Logger.Infof("Sended event %v", event)

}
