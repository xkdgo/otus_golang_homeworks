package internalhttp

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/helper"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
)

func TestCalendarHandler(t *testing.T) {
	storage, err := helper.InitStorage("in-memory", "")
	require.NoError(t, err)
	logg := logger.New("DEBUG")
	app := app.New(logg, storage)
	handler := NewCalendarHandler(app, logg)
	type args struct {
		method string
		uri    string
		body   io.Reader
	}
	type Test struct {
		name           string
		args           args
		wantCode       int
		expectedAnswer interface{}
	}

	tests := []Test{
		{
			name:           "root calendar",
			args:           args{"GET", "http://calendar", nil},
			wantCode:       http.StatusOK,
			expectedAnswer: "Hello, Calendar",
		},
		{
			name:           "api/v1/calendar/event",
			args:           args{"GET", "http://calendar/event", nil},
			wantCode:       http.StatusOK,
			expectedAnswer: "Hello, This is Event Handler",
		},
		{
			name:           "api/v1/calendar/event/create",
			args:           args{"GET", "http://calendar/event/create", nil},
			wantCode:       http.StatusMethodNotAllowed,
			expectedAnswer: "This method not allowed",
		},
		{
			name:           "api/v1/calendar/event/create",
			args:           args{"POST", "http://calendar/event/create", nil},
			wantCode:       http.StatusUnauthorized,
			expectedAnswer: "Bad User",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			r := httptest.NewRequest(tt.args.method, tt.args.uri, tt.args.body)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			resp := w.Result()
			defer resp.Body.Close()
			require.Equal(t, tt.wantCode, resp.StatusCode)
			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tt.expectedAnswer, string(body))
		})
	}

	tests = []Test{
		{
			name:           "api/v1/calendar/event/create",
			args:           args{"POST", "http://calendar/event/create", nil},
			wantCode:       http.StatusBadRequest,
			expectedAnswer: "failed to parse create event\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			r := httptest.NewRequest(tt.args.method, tt.args.uri, tt.args.body)
			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextUserKey, fakeUserID1)
			r = r.WithContext(ctx)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)
			resp := w.Result()
			defer resp.Body.Close()
			require.Equal(t, tt.wantCode, resp.StatusCode)
			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Equal(t, tt.expectedAnswer, string(body))
		})
	}
}
