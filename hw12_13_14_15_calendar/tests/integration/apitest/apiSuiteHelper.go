package apitest

import (
	"fmt"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/tests/integration/generated/openapicli"
)

func (s *APISuiteActions) testEventFromSelfTime() openapicli.EventTemplate {
	t := openapicli.EventTemplate{
		Id:            "",
		Title:         faker.Sentence(),
		Datetimestart: s.eventDate.Format(time.RFC822Z),
		Duration:      fmt.Sprint(time.Duration(24 * time.Hour)),
		Alarmtime:     s.eventDate.Format(time.RFC822Z),
	}
	s.T().Logf("test event: %+v", t)
	return t
}

func (s *APISuiteActions) setNextDay() {
	s.mx.Lock()
	s.eventDate = s.eventDate.Add(24 * time.Hour)
	s.mx.Unlock()
}

func (s *APISuiteActions) ForwardCursorToNextMonday() {
	s.mx.Lock()
	// iterate forward to Monday
	for s.cursorTime.Weekday() != time.Monday {
		s.cursorTime = s.cursorTime.AddDate(0, 0, 1)
	}
	s.mx.Unlock()
}

func (s *APISuiteActions) CountDaysBeforeBeginNextMonthFromCursor() int {
	s.mx.Lock()
	tempDate := s.cursorTime
	s.mx.Unlock()
	daysBeforeNextMonth := 0
	// iterate to begin next month
	if s.cursorTime.Day() == 1 {
		daysBeforeNextMonth++
		tempDate = tempDate.AddDate(0, 0, 1)
	}
	for tempDate.Day() != 1 {
		tempDate = tempDate.AddDate(0, 0, 1)
		daysBeforeNextMonth++
	}
	return daysBeforeNextMonth
}
