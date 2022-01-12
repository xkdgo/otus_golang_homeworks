package sqlstorage

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage/utilstorage"
)

var (
	timelayout        = "2006-01-02"
	timelayoutWithMin = "02 Jan 06 15:04 -0700"
	userIDFirst       = utilstorage.GenerateUUID()
	userIDSecond      = utilstorage.GenerateUUID()
)

func TestStorageGetEvents(t *testing.T) {
	if os.Getenv("CALENDAR_DB") != "pgx" {
		t.Skip("skipping test; $CALENDAR_DB not pgx")
	}
	var testdata storage.Event
	t.Run("test storage get events by day, week, month", func(t *testing.T) {
		var (
			daysInJan = 31
			daysInFeb = 28
		)
		sqlst := New()
		err := sqlst.Connect(context.Background(), "postgres://hwuser:hwpasswd@0.0.0.0:5432/calendar?sslmode=disable")
		require.NoError(t, err)
		defer sqlst.Close()

		err = sqlst.ResetAllData()
		require.NoError(t, err)
		for i := 1; i <= daysInJan; i++ {
			err := faker.FakeData(&testdata)
			require.NoError(t, err)
			testdata.ID = utilstorage.GenerateUUID()
			testdata.UserID = userIDFirst
			testdata.Duration = testdata.Duration * time.Hour

			testdata.DateTimeStart, err = time.Parse(timelayoutWithMin, fmt.Sprintf("%02d Jan 22 12:15 +0500", i))
			require.NoError(t, err)
			_, err = sqlst.CreateEvent(testdata)
			require.NoError(t, err)
		}

		for i := 1; i <= daysInFeb; i++ {
			err := faker.FakeData(&testdata)
			require.NoError(t, err)
			testdata.ID = utilstorage.GenerateUUID()
			testdata.UserID = userIDSecond
			testdata.DateTimeStart, err = time.Parse(timelayoutWithMin, fmt.Sprintf("%02d Feb 22 12:15 +0500", i))
			require.NoError(t, err)
			_, err = sqlst.CreateEvent(testdata)
			require.NoError(t, err)
		}

		testTime, err := time.Parse(timelayout, "2021-01-02")
		require.NoError(t, err)
		scheduledEventsForUser, err := sqlst.ListEventsOnDay(userIDFirst, testTime)
		require.NoError(t, err)
		require.Equal(t, 0, len(scheduledEventsForUser))
		testTime, err = time.Parse(timelayout, "2022-01-02")
		require.NoError(t, err)
		scheduledEventsForUser, err = sqlst.ListEventsOnDay(userIDFirst, testTime)
		require.NoError(t, err)
		require.Equal(t, 1, len(scheduledEventsForUser))
		testTime, err = time.Parse(timelayout, "2022-01-25")
		require.NoError(t, err)
		scheduledEventsForUser, err = sqlst.ListEventsOnCurrentWeek(userIDFirst, testTime)
		require.NoError(t, err)
		require.Equal(t, 6, len(scheduledEventsForUser))
		testTime, err = time.Parse(timelayout, "2021-12-27")
		require.NoError(t, err)
		scheduledEventsForUser, err = sqlst.ListEventsOnCurrentWeek(userIDFirst, testTime)
		require.NoError(t, err)
		require.Equal(t, 2, len(scheduledEventsForUser))
		testTime, err = time.Parse(timelayout, "2022-02-02")
		require.NoError(t, err)
		scheduledEventsForUser, err = sqlst.ListEventsOnCurrentMonth(userIDSecond, testTime)
		require.NoError(t, err)
		require.Equal(t, 27, len(scheduledEventsForUser))
		_, err = sqlst.ListEventsOnCurrentMonth(utilstorage.GenerateUUID(), testTime)
		require.Error(t, err)
		require.Equal(t, "unknown user", err.Error())
	})
}
