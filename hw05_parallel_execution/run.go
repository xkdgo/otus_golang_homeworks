package hw05parallelexecution

import (
	"errors"
	"fmt"

	// "fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var wg sync.WaitGroup

type Task func() error

func produce(taskCh chan Task,
	taskcounter, errorcounter, isErrorEnd *int32,
	lentasks, m int32) {
	defer wg.Done()
	for {
		select {
		case task := <-taskCh:
			switch {
			case task == nil:
				return
			case atomic.LoadInt32(errorcounter) >= m && m > 0:
				atomic.AddInt32(isErrorEnd, 1)
				return
			case atomic.LoadInt32(taskcounter) >= lentasks:
				return
			default:
				err := task()
				if err != nil {
					atomic.AddInt32(errorcounter, 1)

				}
				atomic.AddInt32(taskcounter, 1)
				if atomic.LoadInt32(taskcounter) >= lentasks {
					fmt.Println(atomic.LoadInt32(taskcounter))
					return
				}

			}

		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.
	var taskcounter, errorcounter, isErrorEnd int32
	taskCh := make(chan Task, len(tasks))

	for i := 0; i < n; i++ {
		wg.Add(1)
		go produce(taskCh,
			&taskcounter,
			&errorcounter,
			&isErrorEnd,
			int32(len(tasks)),
			int32(m),
		)
	}
	for _, task := range tasks {
		task := task
		select {
		case taskCh <- task:
		}
	}
	close(taskCh)
	wg.Wait()
	if isErrorEnd >= 1 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
