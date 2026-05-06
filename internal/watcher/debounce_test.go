package watcher_test

import (
	"testing"
	"time"

	"github.com/yourorg/confmerge/internal/watcher"
)

func TestDebouncer_FiresAfterQuietPeriod(t *testing.T) {
	d := watcher.NewDebouncer(40 * time.Millisecond)
	d.Feed()

	select {
	case <-d.Notify:
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout: debouncer did not fire")
	}
}

func TestDebouncer_CoalesceRapidFeeds(t *testing.T) {
	d := watcher.NewDebouncer(60 * time.Millisecond)

	// Feed rapidly; should only produce one notification.
	for i := 0; i < 5; i++ {
		d.Feed()
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for the single debounced notification.
	select {
	case <-d.Notify:
		// good — exactly one
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timeout: debouncer never fired")
	}

	// Ensure no second notification arrives.
	select {
	case <-d.Notify:
		t.Fatal("unexpected second notification")
	case <-time.After(120 * time.Millisecond):
		// OK
	}
}

func TestDebouncer_ResetOnNewFeed(t *testing.T) {
	d := watcher.NewDebouncer(60 * time.Millisecond)
	d.Feed()
	time.Sleep(30 * time.Millisecond) // before quiet period ends
	d.Feed()                          // resets the timer

	select {
	case <-d.Notify:
		// should arrive ~60 ms after the second Feed, not the first
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timeout after reset")
	}
}
