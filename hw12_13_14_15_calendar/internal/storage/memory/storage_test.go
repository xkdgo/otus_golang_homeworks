package memorystorage

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
)

func TestStorage(t *testing.T) {
	memstorage := New()
	t.Run("test create event", func(t *testing.T) {
		testdata := storage.Event{
			Title: "some test event",
		}

		id, err := memstorage.CreateEvent(testdata)
		require.NoError(t, err)
		require.Equal(t, "1", id)
		fmt.Printf("%#v\n", memstorage.data[1])

		anotherData, err := storage.NewEvent("test title",
			"some test desc",
			time.Now().Add(time.Hour*1),
			time.Hour,
			time.Minute)
		require.NoError(t, err)
		id, err = memstorage.CreateEvent(anotherData)
		require.NoError(t, err)
		require.Equal(t, "2", id)
		fmt.Printf("%#v\n", memstorage.data[2])
	})

	t.Run("test overlap id", func(t *testing.T) {
		testdata := storage.Event{
			ID:    "1",
			Title: "overlap test event",
		}
		id, err := memstorage.CreateEvent(testdata)
		require.Error(t, err)
		require.Equal(t, "overlap id of event", err.Error())
		require.Equal(t, "", id)
	})
}
