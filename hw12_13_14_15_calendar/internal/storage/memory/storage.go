package memorystorage

import (
	"sync"
	"time"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage/utilstorage"
)

const layout = "2006-01-02"

type schedule map[time.Time]*storage.Event

type Storage struct {
	mu           sync.RWMutex
	userSchedule map[string]schedule
	data         map[string]*storage.Event
}

func New() *Storage {
	s := &Storage{}
	s.data = make(map[string]*storage.Event)
	s.userSchedule = make(map[string]schedule)
	return s
}

func (s *Storage) ResetAllData() {
	s.data = make(map[string]*storage.Event)
	s.userSchedule = make(map[string]schedule)
}

func (s *Storage) CreateEvent(ev storage.Event) (id string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ev.ID == "" {
		ev.ID = utilstorage.GenerateUUID()
	}
	if ev.UserID == "" {
		return "", storage.ErrEmptyUserIDField
	}
	if _, ok := s.data[ev.ID]; ok {
		return "", storage.ErrOverlapID
	}
	if _, ok := s.userSchedule[ev.UserID]; !ok {
		s.userSchedule[ev.UserID] = make(schedule)
	}
	if _, ok := s.userSchedule[ev.UserID][ev.DateTimeStart]; ok {
		return "", storage.ErrTimeIsBusy
	}
	s.userSchedule[ev.UserID][ev.DateTimeStart] = &ev
	s.data[ev.ID] = &ev
	return s.data[ev.ID].ID, nil
}

func (s *Storage) UpdateEvent(id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ev := event
	if event.ID == "" {
		event.ID = id
	}
	if id != event.ID {
		return storage.ErrMismatchedID
	}
	if _, ok := s.data[event.ID]; !ok {
		return storage.ErrEventIDNotFound
	}
	s.userSchedule[ev.UserID][ev.DateTimeStart] = &ev
	s.data[ev.ID] = &ev
	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event, ok := s.data[id]; ok {
		delete(s.userSchedule[event.UserID], event.DateTimeStart)
	}
	delete(s.data, id)
	return nil
}

func (s *Storage) GetEventsOnDay(userID string, startDate string) (events []*storage.Event, err error) {
	return s.GetEventsByDuration(userID, startDate, 24*time.Hour)
}

func (s *Storage) GetEventsOnCurrentWeek(userID string, startDate string) (events []*storage.Event, err error) {
	dateTime, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, err
	}
	year, weekNumber := dateTime.ISOWeek()
	startYear := time.Date(year, 1, 1, 0, 0, 0, 0, dateTime.Location())
	startOfGivenWeek := startYear.AddDate(0, 0, (weekNumber-1)*7)
	endOfGivenWeek := startOfGivenWeek.AddDate(0, 0, 7)
	durationOfCurrentWeek := endOfGivenWeek.Sub(dateTime)
	return s.GetEventsByDuration(userID, startDate, durationOfCurrentWeek)
}

func (s *Storage) GetEventsOnCurrentMonth(userID string, startDate string) (events []*storage.Event, err error) {
	dateTime, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, err
	}
	year := dateTime.Year()
	month := dateTime.Month()
	startMonth := time.Date(year, month, 1, 0, 0, 0, 0, dateTime.Location())
	endOfGivenMonth := startMonth.AddDate(0, 1, -1)
	durationOfCurrentMonth := endOfGivenMonth.Sub(dateTime)
	return s.GetEventsByDuration(userID, startDate, durationOfCurrentMonth)
}

func (s *Storage) GetEventsByDuration(
	userID string,
	startDate string,
	duration time.Duration,
) (events []*storage.Event, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	userScheduleMap, ok := s.userSchedule[userID]
	if !ok {
		return nil, storage.ErrUnkownUserID
	}
	dateTime, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, err
	}
	dateTimeLater := dateTime.Add(duration)
	for currentDate, ev := range userScheduleMap {
		if currentDate.Before(dateTime) || currentDate.After(dateTimeLater) {
			continue
		} else {
			events = append(events, ev)
		}
	}
	return events, nil
}
