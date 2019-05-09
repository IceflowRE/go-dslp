package utils

import (
	"container/ring"
	"sync"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type CircularBuffer struct {
	Ring   *ring.Ring
	length int
	lock   sync.RWMutex
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		ring.New(size),
		0,
		sync.RWMutex{},
	}
}

// Buffer size.
func (buf *CircularBuffer) Len() int {
	return buf.Ring.Len()
}

// Number of items added.
func (buf *CircularBuffer) Size() int {
	buf.lock.RLock()
	defer buf.lock.RUnlock()
	return buf.length
}

func (buf *CircularBuffer) Current() interface{} {
	buf.lock.RLock()
	defer buf.lock.RUnlock()
	if buf.length == 0 {
		return nil
	}
	return buf.Ring.Value
}

func (buf *CircularBuffer) Add(value interface{}) {
	buf.lock.Lock()
	defer buf.lock.Unlock()
	if buf.length != 0 {
		buf.Ring = buf.Ring.Next()
	}
	buf.Ring.Value = value
	buf.length = min(buf.Ring.Len(), buf.length+1)
}

func (buf *CircularBuffer) Remove() interface{} {
	buf.lock.Lock()
	defer buf.lock.Unlock()
	if buf.length == 0 {
		return nil
	}
	val := buf.Ring.Value
	buf.Ring.Value = nil
	buf.Ring = buf.Ring.Prev()
	buf.length = buf.length - 1
	return val
}
