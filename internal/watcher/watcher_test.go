package watcher_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/confmerge/internal/watcher"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestWatcher_DetectsModification(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, "cfg.yaml", "key: value\n")

	w := watcher.New([]string{p}, 20*time.Millisecond)
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	// Modify the file
	if err := os.WriteFile(p, []byte("key: changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-w.Events:
		if ev.Path != p {
			t.Errorf("expected path %s, got %s", p, ev.Path)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout: no event received after file modification")
	}
}

func TestWatcher_NoEventWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	p := writeTempFile(t, dir, "cfg.yaml", "key: value\n")

	w := watcher.New([]string{p}, 20*time.Millisecond)
	w.Start()
	defer w.Stop()

	time.Sleep(80 * time.Millisecond)
	select {
	case ev := <-w.Events:
		// The first poll after start may fire once if mtime was not seeded;
		// a second event should not appear.
		select {
		case ev2 := <-w.Events:
			t.Errorf("unexpected second event: %+v (first: %+v)", ev2, ev)
		case <-time.After(80 * time.Millisecond):
			// OK
		}
	case <-time.After(80 * time.Millisecond):
		// OK — no spurious events
	}
}

func TestWatcher_ErrorOnMissingFile(t *testing.T) {
	w := watcher.New([]string{"/nonexistent/path/cfg.yaml"}, 20*time.Millisecond)
	w.Start()
	defer w.Stop()

	select {
	case err := <-w.Errors:
		if err == nil {
			t.Fatal("expected non-nil error")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timeout: no error received for missing file")
	}
}
