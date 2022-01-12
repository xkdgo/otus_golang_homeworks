package sqlstorage

import (
	"context"
	"os"
	"testing"

	"sync"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage/utilstorage"
)

func TestStorageConcurency(t *testing.T) { //nolint:gocognit
	if os.Getenv("CALENDAR_DB") != "pgx" {
		t.Skip("skipping test; $CALENDAR_DB not pgx")
	}
	sqlst := New()
	err := sqlst.Connect(context.Background(), "postgres://hwuser:hwpasswd@0.0.0.0:5432/calendar?sslmode=disable")
	require.NoError(t, err)
	defer sqlst.Close()
	err = sqlst.ResetAllData()
	require.NoError(t, err)

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

		sqlst.ResetAllData()
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
					_, err := sqlst.CreateEvent(ev)
					require.NoError(t, err)
				}
			}()
		}
		for _, ev := range testDataSlice {
			eventChan <- ev
		}
		close(eventChan)
		wg.Wait()
		var datacount int64
		rows, err := sqlst.db.Query(`SELECT COUNT(*) AS usercount
		FROM public.events`)
		require.NoError(t, err)

		for rows.Next() {
			err := rows.Scan(&datacount)
			require.NoError(t, err)
		}
		rows.Close()
		require.Equal(t, testDataLen, int(datacount))
		for i := 0; i < workersCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for ev := range updateChan {
					ev.Title = "updated title"
					err := sqlst.UpdateEvent(ev.ID, ev)
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
				query := "select title from public.events where id = $1"
				row := sqlst.db.QueryRow(query, ev.ID)
				var (
					title string
				)
				err := row.Scan(&title)
				require.NoError(t, err)
				require.Equal(t, "updated title", title)
			}
		}

		// delete all updated.
		for i := 0; i < workersCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for ev := range deleteChan {
					err := sqlst.DeleteEvent(ev.ID)
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
		rows, err = sqlst.db.Query(`SELECT COUNT(*) AS usercount
		FROM public.events`)
		require.NoError(t, err)

		for rows.Next() {
			err := rows.Scan(&datacount)
			require.NoError(t, err)
		}
		rows.Close()
		require.Equal(t, testDataLen-updateCounter, int(datacount))
	})
}
