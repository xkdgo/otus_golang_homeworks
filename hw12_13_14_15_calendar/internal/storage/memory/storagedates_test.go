package memorystorage

import (
	"fmt"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

var (
	timelayout        = "2006-01-02"
	timelayoutWithMin = "02 Jan 06 15:04 -0700"
)

func TestStorageGetEvents(t *testing.T) {
	var testdata storage.Event
	t.Run("test storage get events by day, week, month", func(t *testing.T) {
		var (
			daysInJan = 31
			daysInFeb = 28
		)
		memst := New()
		require.Equal(t, 0, len(memst.data))

		for i := 1; i <= daysInJan; i++ {
			err := faker.FakeData(&testdata)
			require.NoError(t, err)
			testdata.UserID = "1"
			testdata.DateTimeStart, err = time.Parse(timelayoutWithMin, fmt.Sprintf("%02d Jan 22 12:15 +0500", i))
			require.NoError(t, err)
			testdata.AlarmTime, err = time.Parse(timelayoutWithMin, fmt.Sprintf("%02d Jan 22 12:00 +0500", i))
			require.NoError(t, err)
			_, err = memst.CreateEvent(testdata)
			require.NoError(t, err)
		}

		require.Equal(t, daysInJan, len(memst.data))
		require.Equal(t, 1, len(memst.userSchedule))

		for i := 1; i <= daysInFeb; i++ {
			err := faker.FakeData(&testdata)
			require.NoError(t, err)
			testdata.UserID = "2"
			testdata.DateTimeStart, err = time.Parse(timelayoutWithMin, fmt.Sprintf("%02d Feb 22 12:15 +0500", i))
			require.NoError(t, err)
			testdata.AlarmTime, err = time.Parse(timelayoutWithMin, fmt.Sprintf("%02d Jan 22 12:00 +0500", i))
			require.NoError(t, err)
			_, err = memst.CreateEvent(testdata)
			require.NoError(t, err)
		}

		require.Equal(t, daysInJan+daysInFeb, len(memst.data))
		require.Equal(t, 2, len(memst.userSchedule))
		testTime, err := time.Parse(timelayout, "2021-01-02")
		require.NoError(t, err)
		scheduledEventsForUser, err := memst.ListEventsOnDay("1", testTime)
		require.NoError(t, err)
		require.Equal(t, 0, len(scheduledEventsForUser))
		testTime, err = time.Parse(timelayout, "2022-01-02")
		require.NoError(t, err)
		scheduledEventsForUser, err = memst.ListEventsOnDay("1", testTime)
		require.NoError(t, err)
		require.Equal(t, 1, len(scheduledEventsForUser))
		testTime, err = time.Parse(timelayout, "2022-01-25")
		require.NoError(t, err)
		scheduledEventsForUser, err = memst.ListEventsOnCurrentWeek("1", testTime)
		require.NoError(t, err)
		require.Equal(t, 6, len(scheduledEventsForUser))
		testTime, err = time.Parse(timelayout, "2021-12-27")
		require.NoError(t, err)
		scheduledEventsForUser, err = memst.ListEventsOnCurrentWeek("1", testTime)
		require.NoError(t, err)
		require.Equal(t, 2, len(scheduledEventsForUser))
		testTime, err = time.Parse(timelayout, "2022-02-02")
		require.NoError(t, err)
		scheduledEventsForUser, err = memst.ListEventsOnCurrentMonth("2", testTime)
		require.NoError(t, err)
		require.Equal(t, 27, len(scheduledEventsForUser))
		_, err = memst.ListEventsOnCurrentMonth("3", testTime)
		require.Error(t, err)
		require.Equal(t, "unknown user", err.Error())
	})
}
