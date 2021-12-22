package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	wg                     sync.WaitGroup
)

type Task func() error

func produce(taskCh chan Task,
	taskCounter, errorCounter, isErrorEnd *int32,
	lentasks, m int32) {
	defer wg.Done()
	for task := range taskCh {
		switch {
		case m == 0 && atomic.LoadInt32(errorCounter) > m:
			atomic.AddInt32(isErrorEnd, 1)
			return
		case atomic.LoadInt32(errorCounter) >= m && m > 0:
			atomic.AddInt32(isErrorEnd, 1)
			return
		case atomic.LoadInt32(taskCounter) >= lentasks:
			return
		default:
			err := task()
			if err != nil {
				atomic.AddInt32(errorCounter, 1)
			}
			atomic.AddInt32(taskCounter, 1)
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}
	var taskCounter, errorCounter, isErrorEnd int32
	taskCh := make(chan Task, len(tasks))

	for i := 0; i < n; i++ {
		wg.Add(1)
		go produce(taskCh,
			&taskCounter,
			&errorCounter,
			&isErrorEnd,
			int32(len(tasks)),
			int32(m),
		)
	}
	for _, task := range tasks {
		task := task
		taskCh <- task // put all tasks without block and finish
	}
	// send nil to close goroutine if success
	close(taskCh)
	wg.Wait()
	if isErrorEnd >= 1 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
