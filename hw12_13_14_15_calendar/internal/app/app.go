package app

import (
	"context"
	"time"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	storage storage.Storage
}

type Logger interface { // TODO
	Log(msg ...interface{})
	Info(msg ...interface{})
	Infof(format string, msg ...interface{})
}

func New(logger Logger, storage storage.Storage) *App {
	return &App{Logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string, userID string, DateTimeStart time.Time, Duration time.Duration, AlarmTime time.Duration) (createdID string, err error) {
	return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(id)
}

// TODO
