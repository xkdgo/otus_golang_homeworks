package memorystorage

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	// TODO
	mu     sync.RWMutex
	nextID int
	data   map[int]*storage.Event
}

func New() *Storage {
	s := &Storage{}
	s.nextID = 1
	s.data = make(map[int]*storage.Event)
	return s
}

func (s *Storage) CreateEvent(ev storage.Event) (id string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ev.ID == "" {
		ev.ID = strconv.Itoa(s.nextID)
		s.nextID++
	}
	idInt, err := strconv.Atoi(ev.ID)
	if err != nil {
		return "", err
	}
	if _, ok := s.data[idInt]; ok {
		return "", storage.ErrOverlapID
	}
	s.data[idInt] = &ev
	return s.data[idInt].ID, nil
}

func (s *Storage) UpdateEvent(id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Errorf("UpdateEvent not implemented")
}

func (s *Storage) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Errorf("DeleteEvent not implemented")
}
