package valueholder

import (
	"sync"
	"time"
)

type ValueHolder struct {
	data       []byte
	lastUpdate time.Time
	lastZero   time.Time
	zero       string
	lock       sync.Locker
}

func NewValueHolder(zero string) *ValueHolder {
	return &ValueHolder{
		data:       make([]byte, 0),
		lock:       new(sync.Mutex),
		lastUpdate: time.Now(),
		lastZero:   time.Now(),
		zero:       zero,
	}
}

func (v *ValueHolder) String() string {
	v.lock.Lock()
	defer v.lock.Unlock()
	return string(v.data)
}

func (v *ValueHolder) ToMap() map[string]any {
	return map[string]any{
		"value":       v.String(),
		"last_update": v.lastUpdate.Format(time.RFC3339),
		"last_zero":   v.lastZero.Format(time.RFC3339),
	}
}

func (v *ValueHolder) Write(p []byte) (n int, err error) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.data = p
	v.lastUpdate = time.Now()
	if string(v.data) == v.zero {
		v.lastZero = v.lastUpdate
	}
	return len(p), nil
}
