package run

import (
	"errors"
	"sync"

	"github.com/davidmz/go-conc/dispose"
)

func It(fn func(onDispose dispose.It) error) error {
	disposers := new(dispose.List)
	defer disposers.Dispose()
	return fn(disposers.Add)
}

func ItVal[T any](fn func(onDispose dispose.It) (T, error)) (T, error) {
	disposers := new(dispose.List)
	defer disposers.Dispose()
	return fn(disposers.Add)
}

// Parallel function just runs in parallel the given funcs and collects
// returning errors. Although it can be used on its own, its main purpose is to
// run inside the run.It and run.ItVal.
func Parallel(funcs ...func() error) error {
	var (
		allErrors = make([]error, 0, len(funcs))
		lk        = new(sync.Mutex)
		wg        = new(sync.WaitGroup)
	)

	wg.Add(len(funcs))
	for _, fn := range funcs {
		fn := fn
		go func() {
			defer wg.Done()
			err := fn()
			lk.Lock()
			allErrors = append(allErrors, err)
			lk.Unlock()
		}()
	}
	wg.Wait()
	return errors.Join(allErrors...)
}
