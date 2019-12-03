package tgo

import (
	"sync"
)

// WaitGroupWrapper WaitGroupWrapper
type WaitGroupWrapper struct {
	sync.WaitGroup
}

// Wrap Wrap
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
