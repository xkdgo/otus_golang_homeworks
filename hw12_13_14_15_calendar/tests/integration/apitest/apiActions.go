package apitest

import (
	"context"
	"net/http"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/tests/integration/generated/openapicli"
)

type APISuiteActions struct {
	suite.Suite
	cli    *openapicli.APIClient
	ctx    context.Context
	apiURL string
	// eventTitle string
	eventDate time.Time
}

func (s *APISuiteActions) Init(apiURL string) {
	apiCfg := openapicli.NewConfiguration()
	s.cli = openapicli.NewAPIClient(apiCfg)
	s.cli.ChangeBasePath(apiURL + "/api/v1")
	s.ctx = context.Background()
	s.apiURL = apiURL
	s.eventDate = time.Now().Add(48 * time.Hour)
}

func (s *APISuiteActions) Client() *openapicli.DefaultApiService {
	return s.cli.DefaultApi
}

func (s *APISuiteActions) CalendarGet() {
	s.T().Helper()
	answ, resp, err := s.cli.DefaultApi.CalendarGet(s.ctx)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("Hello World", answ)
}
