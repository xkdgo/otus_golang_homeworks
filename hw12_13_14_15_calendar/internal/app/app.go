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

func (a *App) GetStorage() storage.Storage {
	return a.storage
}

func (a *App) CreateEvent(ctx context.Context,
	id, title, userID string,
	dateTimeStart time.Time,
	duration, alarmTime time.Duration) (createdID string, err error) {
	return a.storage.CreateEvent(storage.Event{
		ID:            id,
		Title:         title,
		UserID:        userID,
		DateTimeStart: dateTimeStart,
		Duration:      duration,
		AlarmTime:     alarmTime,
	})
}

func (a *App) UpdateEvent(ctx context.Context,
	id, title, userID string,
	dateTimeStart time.Time,
	duration, alarmTime time.Duration) (err error) {
	return a.storage.UpdateEvent(id, storage.Event{
		ID:            id,
		Title:         title,
		UserID:        userID,
		DateTimeStart: dateTimeStart,
		Duration:      duration,
		AlarmTime:     alarmTime,
	})
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(id)
}

func (a *App) ListEventsDay(
	ctx context.Context,
	userID string,
	dateTime time.Time) (events []storage.Event, err error) {
	return a.storage.ListEventsOnDay(userID, dateTime)
}

func (a *App) ListEventsWeek(
	ctx context.Context,
	userID string,
	dateTime time.Time) (events []storage.Event, err error) {
	return a.storage.ListEventsOnCurrentWeek(userID, dateTime)
}

func (a *App) ListEventsMonth(
	ctx context.Context,
	userID string,
	dateTime time.Time) (events []storage.Event, err error) {
	return a.storage.ListEventsOnCurrentMonth(userID, dateTime)
}
