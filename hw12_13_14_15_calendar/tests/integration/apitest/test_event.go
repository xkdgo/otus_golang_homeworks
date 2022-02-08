package apitest

import "sync"

func (s *APISuite) TestCalendarGet() {
	s.CalendarGet()
}

func (s *APISuite) TestCreateEvent() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.CreateEvent()
	}()
	wg.Wait()
	s.CreateEventWithSameDate()
}

func (s *APISuite) TestGetEventsOnCurrentCursorDayWeekMonth() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.CreateEventsOnNdaysForward(60)
		s.ForwardCursorToNextMonday()
	}()
	wg.Wait()
	s.GetEventsOnCurrentCursorDay()
	s.GetEventsOnCurrentCursorWeek()
	s.GetEventsOnCurrentCursorMonth()
}

func (s *APISuite) TestSender() {
	s.CheckMailIsSended()
}
