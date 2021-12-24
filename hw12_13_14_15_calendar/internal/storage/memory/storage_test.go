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

	t.Run("test create overlap id", func(t *testing.T) {
		testdata := storage.Event{
			ID:    "1",
			Title: "overlap test event",
		}
		id, err := memstorage.CreateEvent(testdata)
		require.Error(t, err)
		require.Equal(t, "overlap id of event", err.Error())
		require.Equal(t, "", id)
	})

	t.Run("test create fastforward nextID", func(t *testing.T) {
		memstorage := New()
		for _, EventID := range []string{"1", "2", "3", "5"} {
			testdata := storage.Event{
				ID:    EventID,
				Title: "test event",
			}
			id, err := memstorage.CreateEvent(testdata)
			require.NoError(t, err)
			require.Equal(t, EventID, id)
		}

		testdataAgain := storage.Event{
			Title: "test fastforward event",
		}
		id, err := memstorage.CreateEvent(testdataAgain)
		require.NoError(t, err)
		require.Equal(t, "4", id)
		testdataAgain6 := storage.Event{
			Title: "test fastforward event",
		}
		id, err = memstorage.CreateEvent(testdataAgain6)
		require.NoError(t, err)
		require.Equal(t, "6", id)
	})
}
