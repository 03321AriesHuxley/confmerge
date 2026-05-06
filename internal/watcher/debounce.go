package watcher

import (
	"sync"
	"time"
)

// Debouncer coalesces rapid Events into a single notification after a quiet
// period, preventing redundant pipeline re-runs during bulk file saves.
type Debouncer struct {
	delay  time.Duration
	mu     sync.Mutex
	timer  *time.Timer
	Notify chan struct{}
}

// NewDebouncer creates a Debouncer with the given quiet-period delay.
func NewDebouncer(delay time.Duration) *Debouncer {
	return &Debouncer{
		delay:  delay,
		Notify: make(chan struct{}, 1),
	}
}

// Feed signals the debouncer that an event occurred. The Notify channel
// receives a value only after no further calls arrive within the delay window.
func (d *Debouncer) Feed() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.delay, func() {
		select {
		case d.Notify <- struct{}{}:
		default:
		}
	})
}
