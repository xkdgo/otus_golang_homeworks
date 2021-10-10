package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunFuncsWithErrors(t *testing.T) {
	type args struct {
		tasksWithErrors    int
		tasksWithoutErrors int
		workersCount       int
		maxErrorsCount     int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"negative counter is ignore error",
			args{tasksWithErrors: 1, tasksWithoutErrors: 10, workersCount: 1, maxErrorsCount: -1},
		},

		{
			"negative counter is ignore error",
			args{tasksWithErrors: 100, tasksWithoutErrors: 10, workersCount: 100, maxErrorsCount: -1},
		},
		{
			"negative counter is ignore error",
			args{tasksWithErrors: 100, tasksWithoutErrors: 0, workersCount: 10, maxErrorsCount: -1},
		},
		{
			"error maxcounter is 0",
			args{tasksWithErrors: 7, tasksWithoutErrors: 10, workersCount: 10, maxErrorsCount: 0},
		},
		{
			"error maxcounter is 5",
			args{tasksWithErrors: 8, tasksWithoutErrors: 10, workersCount: 4, maxErrorsCount: 5},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var runTasksCount int32
			tasks := make([]Task, 0, tc.args.tasksWithErrors+tc.args.tasksWithoutErrors)
			for i := 0; i < tc.args.tasksWithErrors; i++ {
				err := fmt.Errorf("error from task %d", i)
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return err
				})
			}

			for i := 0; i < tc.args.tasksWithoutErrors; i++ {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			err := Run(tasks, tc.args.workersCount, tc.args.maxErrorsCount)
			switch {
			case tc.args.maxErrorsCount < 0:
				require.NoError(t, err)
				require.Equal(t, runTasksCount,
					int32(tc.args.tasksWithErrors+tc.args.tasksWithoutErrors), "not all tasks were completed")
			case tc.args.maxErrorsCount == 0:
				require.Error(t, err)
				require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
				require.LessOrEqual(t, runTasksCount, int32(tc.args.workersCount+2), "extra tasks were started")
			case tc.args.maxErrorsCount > 0:
				require.Error(t, err)
				require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
				require.LessOrEqual(t, runTasksCount,
					int32(tc.args.workersCount+tc.args.maxErrorsCount), "extra tasks were started")
			}
		})
	}
}

func TestRunTimeOfExecution(t *testing.T) {
	type args struct {
		tasksWithoutErrors int
		workersCount       int
		maxErrorsCount     int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"negative counter is ignore error",
			args{tasksWithoutErrors: 10, workersCount: 5, maxErrorsCount: -1},
		},
		{
			"negative counter is ignore error",
			args{tasksWithoutErrors: 10, workersCount: 5, maxErrorsCount: -1},
		},
		{
			"negative counter is ignore error",
			args{tasksWithoutErrors: 50, workersCount: 25, maxErrorsCount: -1},
		},
		{
			"error maxcounter is 0",
			args{tasksWithoutErrors: 10, workersCount: 5, maxErrorsCount: 0},
		},
		{
			"error maxcounter is 5",
			args{tasksWithoutErrors: 10, workersCount: 5, maxErrorsCount: 5},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var runTasksCount int32
			tasks := make([]Task, 0, tc.args.tasksWithoutErrors)

			timeOfTask := time.Millisecond * 100
			for i := 0; i < tc.args.tasksWithoutErrors; i++ {
				tasks = append(tasks, func() error {
					time.Sleep(timeOfTask)
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}
			require.Eventually(t, func() bool {
				err := Run(tasks, tc.args.workersCount, tc.args.maxErrorsCount)
				return err == nil
			}, (timeOfTask*time.Duration(tc.args.tasksWithoutErrors))/3, 10*time.Millisecond)
			require.Equal(t, runTasksCount, int32(tc.args.tasksWithoutErrors), "not all tasks were completed")
		})
	}
}
