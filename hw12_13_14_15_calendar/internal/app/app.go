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

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Trace(args ...interface{})
	Tracef(template string, args ...interface{})
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
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
	duration time.Duration,
	alarmTime time.Time) (createdID string, err error) {
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
	duration time.Duration,
	alarmTime time.Time) (err error) {
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
