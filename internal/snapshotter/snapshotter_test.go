package snapshotter_test

import (
	"os"
	"testing"

	"github.com/yourorg/confmerge/internal/snapshotter"
)

func TestSave_And_Load_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	s, err := snapshotter.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	data := map[string]interface{}{"host": "localhost", "port": 8080}
	if err := s.Save("base", data); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := s.Load("base")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap.Name != "base" {
		t.Errorf("expected name 'base', got %q", snap.Name)
	}
	if snap.Data["host"] != "localhost" {
		t.Errorf("expected host 'localhost', got %v", snap.Data["host"])
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := t.TempDir()
	s, _ := snapshotter.New(dir)
	_, err := s.Load("missing")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestList_ReturnsNames(t *testing.T) {
	dir := t.TempDir()
	s, _ := snapshotter.New(dir)

	_ = s.Save("alpha", map[string]interface{}{"a": 1})
	_ = s.Save("beta", map[string]interface{}{"b": 2})

	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(names))
	}
}

func TestList_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	s, _ := snapshotter.New(dir)
	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}

func TestDelete_RemovesSnapshot(t *testing.T) {
	dir := t.TempDir()
	s, _ := snapshotter.New(dir)

	_ = s.Save("temp", map[string]interface{}{"x": true})
	if err := s.Delete("temp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Load("temp")
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestDelete_MissingIsNoop(t *testing.T) {
	dir := t.TempDir()
	s, _ := snapshotter.New(dir)
	if err := s.Delete("nonexistent"); err != nil {
		t.Errorf("expected no error deleting missing snapshot, got %v", err)
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	nested := dir + "/deep/snaps"
	_, err := snapshotter.New(nested)
	if err != nil {
		t.Fatalf("New with nested dir: %v", err)
	}
	if _, err := os.Stat(nested); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
