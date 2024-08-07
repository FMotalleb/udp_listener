// Package valueholder provides a struct called ValueHolder that holds a value, its last update time, its last zero update time,
// and a zero value string. It also has a lock for synchronization (race condition prevention).
package valueholder

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// ValueHolder is a ByteWriter that holds the byte value of the last data written to it.
// It provides methods to read and update the value, as well as to retrieve its last update time and zero update time.
// The lock field is used to synchronize access to the data field and prevent race conditions.
type ValueHolder struct {
	data       []byte
	lastUpdate time.Time
	lastZero   time.Time
	zero       string
	lock       *sync.RWMutex
}

// NewValueHolder creates a new instance of ValueHolder with the provided zero value.
// It initializes the data field with an empty slice, the lastUpdate and lastZero fields with the current time,
// and the zero field with the provided zero value.
func NewValueHolder(zero string) *ValueHolder {
	return &ValueHolder{
		data:       []byte(zero),
		lock:       new(sync.RWMutex),
		lastUpdate: time.Now(),
		lastZero:   time.Now(),
		zero:       zero,
	}
}

// String returns the current value of the ValueHolder as a string.
// It acquires a lock before reading the data field and releases it afterwards.
func (v *ValueHolder) String() string {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return string(v.data)
}

// ToMap returns a map containing the current value, whether it's zero, the last update time,
// and the last zero update time, formatted according to RFC3339.
func (v *ValueHolder) ToMap() map[string]any {
	return map[string]any{
		"value":       v.String(),
		"is_zero":     v.String() == v.zero,
		"last_update": v.lastUpdate.Format(time.RFC3339),
		"last_zero":   v.lastZero.Format(time.RFC3339),
	}
}

// Write updates the current value of the ValueHolder with the provided data slice.
// It acquires a lock before updating the data field, the lastUpdate field, and releases it afterwards.
// If the new value is equal to the zero value, it updates the lastZero field with the current time.
func (v *ValueHolder) Write(p []byte) (n int, err error) {
	v.lock.Lock()
	defer v.lock.Unlock()
	slog.Debug(fmt.Sprintf("rewrite value: `%s` from client", string(p)))
	v.data = p
	v.lastUpdate = time.Now()
	if string(v.data) == v.zero {
		v.lastZero = v.lastUpdate
	}
	return len(p), nil
}
