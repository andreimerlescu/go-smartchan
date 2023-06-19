package go_smartchan

import (
	"errors"
	`fmt`
	"sync"
	`sync/atomic`
)

type SmartChan struct {
	ch     chan interface{}
	mu     sync.RWMutex
	closed atomic.Bool
	cnt    atomic.Int64
}

func NewSmartChan(i int) *SmartChan {
	sc := SmartChan{
		ch: make(chan interface{}, i),
	}
	sc.mu.Lock()
	sc.cnt.Store(0)
	sc.closed.Store(false)
	sc.mu.Unlock()
	return &sc
}

func (sc *SmartChan) Write(thing interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred while writing to the channel: %v", r)
		}
	}()

	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.closed.Load() {
		return errors.New("cannot write to closed channel")
	}
	sc.ch <- thing
	sc.cnt.Add(1)
	return nil
}

func (sc *SmartChan) Read() (interface{}, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	if sc.closed.Load() {
		return nil, errors.New("cannot read from closed channel")
	}
	select {
	case val, ok := <-sc.ch:
		if !ok {
			return nil, errors.New("cannot read from closed or empty channel")
		}
		sc.cnt.Add(-1)
		return val, nil
	default:
		return nil, errors.New("no data available to read")
	}
}

func (sc *SmartChan) Count() int64 {
	return sc.cnt.Load()
}

func (sc *SmartChan) Close() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if !sc.closed.Load() {
		close(sc.ch)
		sc.closed.Store(true)
	}
}

func (sc *SmartChan) Chan() chan interface{} {
	return sc.ch
}

func (sc *SmartChan) CanWrite() bool {
	return !sc.closed.Load()
}
