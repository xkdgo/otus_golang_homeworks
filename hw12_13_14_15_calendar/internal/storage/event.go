package storage

import "time"

type Event struct {
	ID            string
	Title         string
	UserID        string
	DateTimeStart time.Time
	Duration      time.Duration
	AlarmTime     time.Duration
}

func NewEvent(
	title string,
	description string,
	start time.Time,
	duration time.Duration,
	alarm time.Duration) (Event, error) {
	switch {
	case title == "":
		return Event{}, ErrTitle
	case start.Before(time.Now().Add(time.Minute * 15)):
		return Event{}, ErrWithPlannedTime
	}
	e := Event{Title: title, DateTimeStart: start, Duration: duration, AlarmTime: alarm}
	return e, nil
}
