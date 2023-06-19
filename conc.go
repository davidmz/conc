package conc

import (
	"errors"
	"sync"

	"github.com/davidmz/go-conc/dispose"
)

// Task is a function that does some job. It can create resources that need to
// be disposed at the end of process. Disposers of these resources should be
// passed to the 'onDispose' function.
type Task func(onDispose dispose.It) error

// Tasks creates a task from a several other tasks. This tasks will be run in
// parallel.
func Tasks(tasks ...Task) Task {
	return func(onDispose dispose.It) error {
		var (
			allErrors = make([]error, 0, len(tasks))
			lk        = new(sync.Mutex)
			wg        = new(sync.WaitGroup)
		)

		wg.Add(len(tasks))
		for _, task := range tasks {
			task := task
			go func() {
				defer wg.Done()
				err := task(onDispose)
				lk.Lock()
				allErrors = append(allErrors, err)
				lk.Unlock()
			}()
		}
		wg.Wait()
		return errors.Join(allErrors...)
	}
}

// Run executes (in parallel) the given tasks. It calls all collected disposers
// before return.
func Run(task Task, moreTasks ...Task) error {
	disposers := new(dispose.List)
	defer disposers.Dispose()
	if len(moreTasks) == 0 {
		return task(disposers.Add)
	}
	ts := append([]Task{task}, moreTasks...)
	return Tasks(ts...)(disposers.Add)
}
