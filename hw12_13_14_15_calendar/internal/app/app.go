package app

import (
	"context"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

type Logger interface { // TODO

}

type Storage interface {
	CreateEvent(storage.Event) (id string, err error)
	UpdateEvent(id string, event storage.Event) error
	DeleteEvent(id string) error
	// ListEventsDay(time.Time) ([]storage.Event, error)
	// ListEventsWeek(time.Time) ([]storage.Event, error)
	// ListEventsMonth(time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
