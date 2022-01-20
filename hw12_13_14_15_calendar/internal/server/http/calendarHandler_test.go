package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/helper"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/plugins/logger/zap"
)

const (
	fakeUserID1 = "e5446547-ab14-482f-ab72-791079690665"
	fakeUserID2 = "933327d2-3b0b-4688-befd-56da81456859"
)

func TestCalendarHandler(t *testing.T) { //nolint:funlen
	storage, err := helper.InitStorage("in-memory", "")
	require.NoError(t, err)
	pluginlogger, err := zap.NewLogger()
	require.NoError(t, err)
	logg := logger.New("DEBUG", pluginlogger)
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
			"root calendar", args{"GET", "http://calendar", nil}, http.StatusOK, "Hello, Calendar",
		},
		{
			"api/v1/calendar/event", args{"GET", "http://calendar/event", nil}, http.StatusOK, "Hello, This is Event Handler",
		},
		{
			"api/v1/calendar/event/create | Method not allowed",
			args{"GET", "http://calendar/event/create", nil},
			http.StatusMethodNotAllowed,
			"This method not allowed",
		},
		{
			"api/v1/calendar/event/create | Unathorized",
			args{"POST", "http://calendar/event/create", nil},
			http.StatusUnauthorized,
			"Bad User",
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

	eventJSONString := `{"id":"",
	"title":"Test Event",
	"datetimestart":"01 Feb 22 12:15 +0500",
	"duration":"26h20m0s",
	"alarmtime":"01 Feb 22 12:00 +0500"}`
	rawIn := json.RawMessage(eventJSONString)
	bytesEncoded, err := rawIn.MarshalJSON()
	require.NoError(t, err)
	buf := bytes.Buffer{}
	buf.Write(bytesEncoded)

	tests = []Test{
		{
			name:           "api/v1/calendar/event/create | bad request",
			args:           args{"POST", "http://calendar/event/create", nil},
			wantCode:       http.StatusBadRequest,
			expectedAnswer: "failed to parse create event\n",
		},
		{
			name:           "api/v1/calendar/event/create | without eventID",
			args:           args{"POST", "http://calendar/event/create", &buf},
			wantCode:       http.StatusOK,
			expectedAnswer: "Hello, your event is created with id",
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
			// require.Equal(t, tt.expectedAnswer, string(body))
			require.Contains(t, string(body), tt.expectedAnswer)
		})
	}

	eventWithID := `{"id":"8abb4997-2dc1-4bf3-b6ca-fe27b12724dd",
	"title":"Test Event",
	"datetimestart":"01 Feb 22 12:16 +0500",
	"duration":"26h20m0s",
	"alarmtime":"01 Feb 22 12:00 +0500"}`
	rawIn = json.RawMessage(eventWithID)
	bytesEncoded, err = rawIn.MarshalJSON()
	require.NoError(t, err)
	buf = bytes.Buffer{}
	buf.Write(bytesEncoded)

	t.Run("create eventWithID", func(t *testing.T) {
		r := httptest.NewRequest("POST", "http://calendar/event/create", &buf)
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserKey, fakeUserID1)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	eventWithID = `{"id":"8abb4997-2dc1-4bf3-b6ca-fe27b12724dd",
	"title":"Test Event Update",
	"datetimestart":"01 Feb 22 12:16 +0500",
	"duration":"26h20m0s",
	"alarmtime":"01 Feb 22 12:00 +0500"}`
	rawIn = json.RawMessage(eventWithID)
	bytesEncoded, err = rawIn.MarshalJSON()
	require.NoError(t, err)
	buf = bytes.Buffer{}
	buf.Write(bytesEncoded)

	t.Run("good update eventWithID", func(t *testing.T) {
		r := httptest.NewRequest("POST", "http://calendar/event/update", &buf)
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserKey, fakeUserID1)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	eventWithID = `{"id":"8abb4997-2dc1-4bf3-b6ca-ffffffffffff",
	"title":"Test Event Update",
	"datetimestart":"01 Feb 22 12:16 +0500",
	"duration":"26h20m0s",
	"alarmtime":"01 Feb 22 12:00 +0500"}`
	rawIn = json.RawMessage(eventWithID)
	bytesEncoded, err = rawIn.MarshalJSON()
	require.NoError(t, err)
	buf = bytes.Buffer{}
	buf.Write(bytesEncoded)

	t.Run("bad update eventWithID", func(t *testing.T) {
		r := httptest.NewRequest("POST", "http://calendar/event/update", &buf)
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserKey, fakeUserID1)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("delete eventWithID", func(t *testing.T) {
		r := httptest.NewRequest("POST", "http://calendar/event/delete/8abb4997-2dc1-4bf3-b6ca-fe27b12724dd", nil)
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserKey, fakeUserID1)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("delete again eventWithID status ok anyway", func(t *testing.T) {
		r := httptest.NewRequest("POST", "http://calendar/event/delete/8abb4997-2dc1-4bf3-b6ca-fe27b12724dd", nil)
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserKey, fakeUserID1)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
