package gobatch

import (
	"sync"
	"time"
)

func NewMemoryBatch(flushMaxSize int, flushMaxWait time.Duration, callback BatchFn, workerSize int, mode Mode) *Batch {
	instance := &Batch{
		maxSize: flushMaxSize,
		maxWait: flushMaxWait,

		items: []interface{}{},
		doFn:  callback,
		mutex: &sync.RWMutex{},

		flushChan: make(chan []interface{}, workerSize),
		Mode:      mode,
	}
	instance.setFlushWorker(workerSize)
	go instance.runFlushByTime()
	return instance
}

func (b *Batch) Insert(data interface{}) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.items = append(b.items, data)

	if b.Mode == FlushBySize || b.Mode == FlushByTimeAndSize {
		if len(b.items) >= b.maxSize {
			b.Flush()
		}
	}

}

func (b *Batch) runFlushByTime() {
	for {
		select {
		case <-time.Tick(b.maxWait):
			if b.Mode == FlushByTime || b.Mode == FlushByTimeAndSize {
				b.mutex.Lock()
				b.Flush()
				b.mutex.Unlock()
			}
		}
	}
}

func (b *Batch) Flush() {
	if len(b.items) <= 0 {
		return
	}

	copiedItems := make([]interface{}, len(b.items))
	for idx, i := range b.items {
		copiedItems[idx] = i
	}
	b.items = b.items[:0]
	b.flushChan <- copiedItems
}
