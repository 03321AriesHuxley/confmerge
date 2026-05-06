package watcher

import (
	"os"
	"path/filepath"
	"time"
)

// Event represents a file change event.
type Event struct {
	Path    string
	ModTime time.Time
}

// Watcher polls a set of files for modifications.
type Watcher struct {
	files    []string
	interval time.Duration
	mtimes   map[string]time.Time
	Events   chan Event
	Errors   chan error
	stop     chan struct{}
}

// New creates a Watcher that polls the given files at the given interval.
func New(files []string, interval time.Duration) *Watcher {
	w := &Watcher{
		files:    files,
		interval: interval,
		mtimes:   make(map[string]time.Time),
		Events:   make(chan Event, 16),
		Errors:   make(chan error, 4),
		stop:     make(chan struct{}),
	}
	for _, f := range files {
		if info, err := os.Stat(filepath.Clean(f)); err == nil {
			w.mtimes[f] = info.ModTime()
		}
	}
	return w
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-w.stop:
				return
			case <-ticker.C:
				w.poll()
			}
		}
	}()
}

// Stop halts the polling goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) poll() {
	for _, f := range w.files {
		info, err := os.Stat(filepath.Clean(f))
		if err != nil {
			w.Errors <- err
			continue
		}
		prev, seen := w.mtimes[f]
		if !seen || info.ModTime().After(prev) {
			w.mtimes[f] = info.ModTime()
			w.Events <- Event{Path: f, ModTime: info.ModTime()}
		}
	}
}
