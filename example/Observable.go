// +build ignore

// +gotemplate

package example

import "sync"

type Observable_T struct {
	mutex     sync.RWMutex
	value     T
	observers map[*Observer_T]struct{}
}

func newObservable_T(value T) *Observable_T {
	return &Observable_T{
		value:     value,
		observers: make(map[*Observer_T]struct{}),
	}
}

func (o *Observable_T) Set(value T) {
	o.mutex.Lock()
	o.value = value
	o.mutex.Unlock()

	o.mutex.RLock()
	for obs, _ := range o.observers {
		go obs.Notify(value)
	}
	o.mutex.RUnlock()
}

func (o *Observable_T) Value() T {
	o.mutex.RLock()
	v := o.value
	o.mutex.RUnlock()
	return v
}

func (o *Observable_T) Observe(callback func(T)) *Observer_T {
	o.mutex.Lock()
	obs := &Observer_T{callback, o}
	o.observers[obs] = struct{}{}
	o.mutex.Unlock()
	return obs
}

type Observer_T struct {
	Notify func(T)

	observable *Observable_T
}

func (o *Observer_T) Close() {
	o.observable.mutex.Lock()
	delete(o.observable.observers, o)
	o.observable.mutex.Unlock()
}

type CircularBuffer_T struct {
	data   []T
	curPos int
}

func newCircularBuffer_T(size int) {
	return CircularBuffer_T{data: make([]T, size)}
}

func (b CircularBuffer_T) Push(value T) {
	b.curPos++
	if b.curPos >= len(b.data) {
		b.curPos = 0
	}
	b.data[b.curPos] = value
}

func (b CircularBuffer_T) At(index uint) T {
	pos := b.curPos - index
	for pos < 0 {
		pos += len(b.data)
	}
	return b.data[pos]
}
