package run

import (
	"errors"
	"sync"

	"github.com/davidmz/go-conc/dispose"
	"github.com/davidmz/go-try"
)

func It(fn func(onDispose dispose.It)) (outErr error) {
	defer try.HandleAs(&outErr)
	TryIt(fn)
	return
}

func ItVal[T any](fn func(onDispose dispose.It) T) (result T, outErr error) {
	defer try.HandleAs(&outErr)
	result = TryItVal(fn)
	return
}

func TryIt(fn func(onDispose dispose.It)) {
	disposers := new(dispose.List)
	defer disposers.Dispose()
	fn(disposers.Add)
}

func TryItVal[T any](fn func(onDispose dispose.It) T) T {
	disposers := new(dispose.List)
	defer disposers.Dispose()
	return fn(disposers.Add)
}

func Parallel(funcs ...func()) {
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
			defer try.Handle(func(err error) {

				lk.Lock()
				allErrors = append(allErrors, err)
				lk.Unlock()
			})
			fn()
		}()
	}
	wg.Wait()
	try.It(errors.Join(allErrors...))
}
