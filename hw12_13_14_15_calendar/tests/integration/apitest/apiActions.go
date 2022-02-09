package apitest

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/tests/integration/generated/openapicli"
)

const (
	fakeUserID1 = "e5446547-ab14-482f-ab72-791079690665"
)

type APISuiteActions struct {
	suite.Suite
	mx     *sync.Mutex
	cli    *openapicli.APIClient
	ctx    context.Context
	apiURL string
	// eventTitle string
	eventDate  time.Time
	cursorTime time.Time
}

func (s *APISuiteActions) Init(apiURL string) {
	s.mx = &sync.Mutex{}
	key := openapicli.APIKey{Key: fakeUserID1}
	auth := make(map[string]openapicli.APIKey)
	auth["UserAuth"] = key
	apiCfg := openapicli.NewConfiguration()
	s.cli = openapicli.NewAPIClient(apiCfg)
	s.cli.GetConfig().Servers[0].URL = apiURL + "/api/v1"
	s.ctx = context.Background()
	s.ctx = context.WithValue(s.ctx, openapicli.ContextAPIKeys, auth)
	s.apiURL = apiURL
	s.eventDate = time.Now().Add(1 * time.Hour)
	s.cursorTime = time.Now()
}

func (s *APISuiteActions) Client() *openapicli.DefaultApiService {
	return s.cli.DefaultApi
}

func (s *APISuiteActions) CalendarGet() {
	s.T().Helper()
	answ, resp, err := s.cli.DefaultApi.CalendarGet(s.ctx).Execute()
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("Hello, Calendar", answ)
}

func (s *APISuiteActions) CreateEvent() {
	s.T().Helper()
	answ, resp, err := s.cli.DefaultApi.CreateEvent(s.ctx).EventTemplate(s.testEventFromSelfTime()).Execute()
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Contains(answ, "Hello, your event is created with id")
}

func (s *APISuiteActions) CreateEventWithSameDate() {
	s.T().Helper()
	_, resp, err := s.cli.DefaultApi.CreateEvent(s.ctx).EventTemplate(s.testEventFromSelfTime()).Execute()
	s.Require().Errorf(err, "catched %s", err)
	buf := new(strings.Builder)
	_, errc := io.Copy(buf, resp.Body)
	answ := buf.String()
	s.Require().NoError(errc)
	s.T().Logf("error: %v, answer: %v", err, answ)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
	s.Require().Equal("failed to create event\n", answ)
}

func (s *APISuiteActions) CreateEventsOnNdaysForward(n int) {
	s.T().Helper()
	if n <= 0 {
		n = 1
	}
	for i := 0; i < n; i++ {
		s.setNextDay()
		s.CreateEvent()
	}
}

func (s *APISuiteActions) GetEventsOnCurrentCursorDay() {
	s.T().Helper()
	date := s.cursorTime.Format("2006-01-02")
	answ, resp, err := s.cli.DefaultApi.GetEventsByDay(s.ctx, date).Execute()
	defer resp.Body.Close()
	s.Require().NoError(err)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.T().Logf("answ: %v", answ)
	s.Require().Equal(1, len(answ))
}

func (s *APISuiteActions) GetEventsOnCurrentCursorWeek() {
	s.T().Helper()
	date := s.cursorTime.Format("2006-01-02")
	answ, resp, err := s.cli.DefaultApi.GetEventsByWeek(s.ctx, date).Execute()
	defer resp.Body.Close()
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.T().Logf("answ: %v", answ)
	s.Require().Equal(7, len(answ))
}

func (s *APISuiteActions) GetEventsOnCurrentCursorMonth() {
	s.T().Helper()
	date := s.cursorTime.Format("2006-01-02")
	answ, resp, err := s.cli.DefaultApi.GetEventsByMonth(s.ctx, date).Execute()
	defer resp.Body.Close()
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.T().Logf("answ: %v", answ)
	s.Require().Equal(s.CountDaysBeforeBeginNextMonthFromCursor(), len(answ))
}

func (s *APISuiteActions) CheckMailIsSended() {
	s.T().Helper()
	timeoutEnv := os.Getenv("TESTS_TIMEOUT")
	timeout, err := time.ParseDuration(timeoutEnv)
	s.Require().NoError(err)
	s.mx.Lock()
	s.eventDate = time.Now().Add(timeout)
	s.mx.Unlock()
	s.CreateEvent()
	toPath := os.Getenv("TESTS_MAIL")
	timer := time.NewTimer(timeout * 2)
	var sendDetected bool
	sended := make(chan struct{})
	go func(ch chan struct{}) {
		fd, err := os.Open(toPath)
		s.Require().NoError(err)
		time.Sleep(timeout)
		scanner := bufio.NewScanner(fd)
		for scanner.Scan() {
			answ := scanner.Text()
			s.T().Logf("reading mail log: %v", answ)
			if strings.Contains(answ, "title") {
				close(ch)
				return
			}
		}
	}(sended)
	select {
	case <-timer.C:
		sendDetected = false
	case <-sended:
		sendDetected = true
	}
	s.Require().Truef(sendDetected, "send not detected during timeout %s", timeout*3)
}
