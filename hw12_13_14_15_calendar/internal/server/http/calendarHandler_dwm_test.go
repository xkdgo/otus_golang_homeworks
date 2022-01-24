package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/helper"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/logger"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/server/http/internal/models"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/plugins/logger/zap"
)

func TestCalendarDayWeekMonthHandler(t *testing.T) {
	storage, err := helper.InitStorage("in-memory", "")
	require.NoError(t, err)
	pluginlogger, err := zap.NewLogger()
	require.NoError(t, err)
	logg := logger.New("DEBUG", pluginlogger)
	app := app.NewAppCalendar(logg, storage)
	handler := NewCalendarHandler(app, logg)
	var buf bytes.Buffer
	const (
		DaysInFeb      = 28
		NumsOfDayTasks = 8
		StartHour      = 8
		daysInWeek     = 7
	)
	for day := 1; day < DaysInFeb+1; day++ {
		for hour := 8; hour < StartHour+NumsOfDayTasks; hour++ {
			event := fmt.Sprintf(`{"id":"",
			"title":"Test Event",
			"datetimestart":"%02d Feb 22 %02d:16 +0500",
			"duration":"1h",
			"alarmtime":"%02d Feb 22 %02d:00 +0500"}`, day, hour, day, hour)
			rawIn := json.RawMessage(event)
			bytesEncoded, err := rawIn.MarshalJSON()
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
		}
	}

	// Test get events per day
	t.Run("Test get events per day", func(t *testing.T) {
		var result []models.Event
		r := httptest.NewRequest("GET", "http://calendar/event/day/2022-02-01", nil)
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserKey, fakeUserID1)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)
		require.Len(t, result, NumsOfDayTasks)
	})

	// Test get events per week
	t.Run("Test get events per week", func(t *testing.T) {
		var result []models.Event
		r := httptest.NewRequest("GET", "http://calendar/event/week/2022-02-01", nil)
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserKey, fakeUserID1)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)
		require.Len(t, result, NumsOfDayTasks*(daysInWeek-1))
	})

	// Test get events per month
	t.Run("Test get events per month", func(t *testing.T) {
		var result []models.Event
		r := httptest.NewRequest("GET", "http://calendar/event/month/2022-02-01", nil)
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserKey, fakeUserID1)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)
		require.Len(t, result, NumsOfDayTasks*DaysInFeb)
	})

	// Test get events per month
	t.Run("Test get events per month with unknown user", func(t *testing.T) {
		r := httptest.NewRequest("GET", "http://calendar/event/month/2022-02-01", nil)
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextUserKey, fakeUserID2)
		r = r.WithContext(ctx)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		resp := w.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
