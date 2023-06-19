package dispose

import "sync"

// Disposer is a function that frees some resource. It will be called at the end
// of whole Run sequence (as a 'defer'). Task can register Disposer using it
// 'onDispose' argument.
type Disposer func()

// It is a function that takes Disposer.
type It func(Disposer)

type List struct {
	sync.Mutex
	disposers []Disposer
}

func (ds *List) Dispose() {
	for _, d := range ds.disposers {
		d()
	}
}

func (ds *List) Add(d Disposer) {
	ds.Lock()
	ds.disposers = append(ds.disposers, d)
	ds.Unlock()
}
