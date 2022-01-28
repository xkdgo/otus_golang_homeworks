package storage

import "time"

type Storage interface {
	Close() error
	CreateEvent(Event) (id string, err error)
	UpdateEvent(id string, event Event) error
	DeleteEvent(id string) error
	ListEventsOnDay(userID string, dateTime time.Time) (events []Event, err error)
	ListEventsOnCurrentWeek(userID string, dateTime time.Time) (events []Event, err error)
	ListEventsOnCurrentMonth(userID string, dateTime time.Time) (events []Event, err error)
	ListEventsToNotify(periodTimeStart time.Time, periodTimeEnd time.Time) (events []Event, err error)
	ListEventsToDelete(ttl time.Duration) (events []Event, err error)
}
