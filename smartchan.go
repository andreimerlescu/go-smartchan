package go_smartchan

import (
	"errors"
	"sync"
)

type SmartChan struct {
	ch     chan interface{}
	mu     sync.RWMutex
	closed bool
	cnt    int64
}

func NewSmartChan(i int) *SmartChan {
	return &SmartChan{
		ch:     make(chan interface{}, i),
		closed: false,
		cnt:    0,
	}
}

func (sc *SmartChan) Write(thing interface{}) error {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	if sc.closed {
		return errors.New("cannot write to closed channel")
	}
	sc.ch <- thing
	sc.cnt++
	return nil
}

func (sc *SmartChan) Read() (interface{}, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	if sc.closed && len(sc.ch) == 0 {
		return nil, errors.New("cannot read from closed channel")
	}
	return <-sc.ch, nil
}

func (sc *SmartChan) Count() int64 {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.cnt
}

func (sc *SmartChan) Close() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if !sc.closed {
		close(sc.ch)
		sc.closed = true
	}
}

func (sc *SmartChan) Chan() chan interface{} {
	return sc.ch
}

func (sc *SmartChan) CanWrite() bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return !sc.closed
}
