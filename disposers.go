package conc

import "sync"

// Disposer is a function that frees some resource. It will be called at the end
// of whole Run sequence (as a 'defer'). Task can register Disposer using it
// 'onDispose' argument.
type Disposer func()

// OnDispose is a function that takes Disposer.
type OnDispose func(Disposer)

type disposersList struct {
	sync.Mutex
	disposers []Disposer
}

func (ds *disposersList) Dispose() {
	for _, d := range ds.disposers {
		d()
	}
}

func (ds *disposersList) Add(d Disposer) {
	ds.Lock()
	ds.disposers = append(ds.disposers, d)
	ds.Unlock()
}
