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
)

const (
	fakeUserID1 = "e5446547-ab14-482f-ab72-791079690665"
	fakeUserID2 = "933327d2-3b0b-4688-befd-56da81456859"
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
			name:           "api/v1/calendar/event/create | Method not allowed",
			args:           args{"GET", "http://calendar/event/create", nil},
			wantCode:       http.StatusMethodNotAllowed,
			expectedAnswer: "This method not allowed",
		},
		{
			name:           "api/v1/calendar/event/create | Unathorized",
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

	// modelEvent := models.Event{
	// 	ID:            "",
	// 	Title:         "Test Event",
	// 	DateTimeStart: "01 Feb 22 12:15",
	// 	Duration:      models.Duration{Duration: time.Hour*2 + time.Hour*24 + time.Minute*20},
	// 	AlarmTime:     models.Duration{Duration: time.Minute * 15},
	// }

	// modelEventJson, err := json.Marshal(modelEvent)
	// fmt.Println(string(modelEventJson))
	// 5d62f3b3-923c-4514-93a3-64c3dd053f0c
	eventJsonString := `{"id":"",
	"title":"Test Event",
	"datetimestart":"01 Feb 22 12:15 +0500",
	"duration":"26h20m0s",
	"alarmtime":"15m0s"}`
	rawIn := json.RawMessage(eventJsonString)
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
	"alarmtime":"15m0s"}`
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
	"alarmtime":"15m0s"}`
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
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	eventWithID = `{"id":"8abb4997-2dc1-4bf3-b6ca-ffffffffffff",
	"title":"Test Event Update",
	"datetimestart":"01 Feb 22 12:16 +0500",
	"duration":"26h20m0s",
	"alarmtime":"15m0s"}`
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
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
