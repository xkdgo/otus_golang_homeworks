//go:build integration
// +build integration

package integration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/tests/integration/apitest"
)

type IntegrationSuite struct {
	apitest.APISuite
}

func (s *IntegrationSuite) SetupTest() {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://127.0.0.1:8080"
	}

	s.Init(apiURL)

	// s.Client().SetEventsTemplate(context.Background(), openapicli.EventsTemplate{})
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
