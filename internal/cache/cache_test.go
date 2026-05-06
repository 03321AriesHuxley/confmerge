package cache_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/confmerge/internal/cache"
)

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "cache")
	c, err := cache.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Cache")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("expected directory %q to be created", dir)
	}
}

func TestPutAndGet_RoundTrip(t *testing.T) {
	c, _ := cache.New(t.TempDir())

	data := map[string]interface{}{"env": "production", "port": 8080}
	if err := c.Put("abc123", data); err != nil {
		t.Fatalf("Put: %v", err)
	}

	entry, err := c.Get("abc123")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if entry == nil {
		t.Fatal("expected non-nil entry")
	}
	if entry.Key != "abc123" {
		t.Errorf("key: got %q, want %q", entry.Key, "abc123")
	}
	if entry.Data["env"] != "production" {
		t.Errorf("data[env]: got %v, want production", entry.Data["env"])
	}
}

func TestGet_MissReturnsNil(t *testing.T) {
	c, _ := cache.New(t.TempDir())

	entry, err := c.Get("nonexistent")
	if err != nil {
		t.Fatalf("Get: unexpected error: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil on cache miss, got %+v", entry)
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c, _ := cache.New(t.TempDir())

	_ = c.Put("xyz", map[string]interface{}{"a": 1})
	if err := c.Invalidate("xyz"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	entry, err := c.Get("xyz")
	if err != nil {
		t.Fatalf("Get after invalidate: %v", err)
	}
	if entry != nil {
		t.Errorf("expected nil after invalidation, got %+v", entry)
	}
}

func TestInvalidate_MissingKeyIsNoop(t *testing.T) {
	c, _ := cache.New(t.TempDir())
	if err := c.Invalidate("does-not-exist"); err != nil {
		t.Errorf("Invalidate missing key should not error: %v", err)
	}
}

func TestKey_DeterministicForSamePaths(t *testing.T) {
	f := filepath.Join(t.TempDir(), "base.yaml")
	_ = os.WriteFile(f, []byte("a: 1"), 0o644)

	k1, err := cache.Key([]string{f})
	if err != nil {
		t.Fatalf("Key: %v", err)
	}
	k2, err := cache.Key([]string{f})
	if err != nil {
		t.Fatalf("Key: %v", err)
	}
	if k1 != k2 {
		t.Errorf("Key not deterministic: %q != %q", k1, k2)
	}
}
