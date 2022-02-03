package apitest

import (
	"context"
	"net/http"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/tests/integration/generated/openapicli"
)

const (
	fakeUserID1 = "e5446547-ab14-482f-ab72-791079690665"
	fakeUserID2 = "933327d2-3b0b-4688-befd-56da81456859"
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
	key := openapicli.APIKey{Key: fakeUserID1}
	auth := make(map[string]openapicli.APIKey)
	auth["UserAuth"] = key
	apiCfg := openapicli.NewConfiguration()
	s.cli = openapicli.NewAPIClient(apiCfg)
	s.cli.GetConfig().Servers[0].URL = apiURL + "/api/v1"
	s.ctx = context.Background()
	s.ctx = context.WithValue(s.ctx, openapicli.ContextAPIKeys, auth)
	s.apiURL = apiURL
	s.eventDate = time.Now().Add(48 * time.Hour)
}

func (s *APISuiteActions) Client() *openapicli.DefaultApiService {
	return s.cli.DefaultApi
}

func (s *APISuiteActions) CalendarGet() {
	s.T().Helper()
	req := s.cli.DefaultApi.CalendarGet(s.ctx)
	answ, resp, err := s.cli.DefaultApi.CalendarGetExecute(req)
	// answ, resp, err := s.cli.DefaultApi.CalendarGet(s.ctx)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("Hello, Calendar", answ)
}
