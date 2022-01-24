package memorystorage

import (
	"fmt"
	"sync"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage/utilstorage"
)

func TestStorage(t *testing.T) {
	memstorage := New()
	var testdata storage.Event
	t.Run("test create event", func(t *testing.T) {
		require.Equal(t, 0, len(memstorage.data))
		err := faker.FakeData(&testdata)
		testdata.ID = ""
		testdata.UserID = userIDFirst
		require.NoError(t, err)
		_, err = memstorage.CreateEvent(testdata)
		require.NoError(t, err)
		require.Equal(t, 1, len(memstorage.data))
		require.Equal(t, 1, len(memstorage.userSchedule))
		fmt.Printf("%#v\n", memstorage.data)
		fmt.Printf("%#v\n", memstorage.userSchedule)

		err = faker.FakeData(&testdata)
		require.NoError(t, err)
		testdata.ID = ""
		testdata.UserID = userIDSecond
		_, err = memstorage.CreateEvent(testdata)
		require.NoError(t, err)
		require.Equal(t, 2, len(memstorage.data))
		require.Equal(t, 2, len(memstorage.userSchedule))
		fmt.Printf("%#v\n", memstorage.data)
		fmt.Printf("%#v\n", memstorage.userSchedule)
	})

	t.Run("test create overlap id", func(t *testing.T) {
		err := faker.FakeData(&testdata)
		require.NoError(t, err)
		testdata.ID = "00112233-4455-6677-8899-aabbccddeeff"
		testdata.Title = "some test event"
		id, err := memstorage.CreateEvent(testdata)
		require.NoError(t, err)
		require.Equal(t, "00112233-4455-6677-8899-aabbccddeeff", id)
		err = faker.FakeData(&testdata)
		require.NoError(t, err)
		testdata.ID = "00112233-4455-6677-8899-aabbccddeeff"
		testdata.Title = "overlap test event"
		id, err = memstorage.CreateEvent(testdata)
		require.Equal(t, "overlap id of event", err.Error())
		require.Equal(t, "", id)
	})
}

func TestStorageConcurency(t *testing.T) { //nolint:gocognit
	memstorage := New()
	var testdata storage.Event
	t.Run("test concurency", func(t *testing.T) {
		var (
			divider      = 3
			testDataLen  = 1000
			workersCount = 10
			wg           sync.WaitGroup
		)
		eventChan := make(chan storage.Event, 10)
		updateChan := make(chan storage.Event, 10)
		deleteChan := make(chan storage.Event, 10)

		memstorage.ResetAllData()
		testDataSlice := make([]storage.Event, 0, testDataLen)
		for i := 0; i < testDataLen; i++ {
			err := faker.FakeData(&testdata)
			require.NoError(t, err)
			testdata.ID = utilstorage.GenerateUUID()
			testdata.UserID = utilstorage.GenerateUUID()
			testDataSlice = append(testDataSlice, testdata)
		}
		for i := 0; i < workersCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for ev := range eventChan {
					_, err := memstorage.CreateEvent(ev)
					require.NoError(t, err)
				}
			}()
		}
		for _, ev := range testDataSlice {
			eventChan <- ev
		}
		close(eventChan)
		wg.Wait()
		require.Equal(t, testDataLen, len(memstorage.data))
		for i := 0; i < workersCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for ev := range updateChan {
					ev.Title = "updated title"
					err := memstorage.UpdateEvent(ev.ID, ev)
					require.NoError(t, err)
				}
			}()
		}
		// update every divider count.
		var updateCounter int
		for index, ev := range testDataSlice {
			if index%divider == 0 {
				updateChan <- ev
				updateCounter++
			}
		}
		close(updateChan)
		wg.Wait()
		for index, ev := range testDataSlice {
			if index%divider == 0 {
				require.Equal(t, "updated title", memstorage.data[ev.ID].Title)
			}
		}

		// delete all updated.
		for i := 0; i < updateCounter; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for ev := range deleteChan {
					err := memstorage.DeleteEvent(ev.ID)
					require.NoError(t, err)
				}
			}()
		}
		for index, ev := range testDataSlice {
			if index%divider == 0 {
				deleteChan <- ev
			}
		}
		close(deleteChan)
		wg.Wait()
		require.Equal(t, testDataLen-updateCounter, len(memstorage.data))
	})
}
