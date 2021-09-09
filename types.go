package gobatch

import (
	"sync"
	"time"
)

type Mode int

const (
	FlushByTimeAndSize Mode = iota
	FlushByTime
	FlushBySize
)

type BatchFn func(workerID int, datas []interface{})
type Batch struct {
	maxSize int
	maxWait time.Duration

	items []interface{}
	doFn  BatchFn
	mutex *sync.RWMutex

	/*notifier channel*/
	flushChan chan []interface{}
	Mode      Mode
}
