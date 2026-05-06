// Package cache provides a simple file-based result cache for confmerge.
// It stores the merged output keyed by a hash of the input file paths and
// their modification times, allowing repeated runs to skip redundant work.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry holds a cached merge result along with metadata.
type Entry struct {
	Key       string                 `json:"key"`
	CreatedAt time.Time              `json:"created_at"`
	Data      map[string]interface{} `json:"data"`
}

// Cache manages on-disk cache entries under a configurable directory.
type Cache struct {
	dir string
}

// New creates a Cache that stores entries under dir.
// The directory is created if it does not already exist.
func New(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cache: create dir %q: %w", dir, err)
	}
	return &Cache{dir: dir}, nil
}

// Key derives a deterministic cache key from a list of file paths by hashing
// each path together with its last-modified time.
func Key(paths []string) (string, error) {
	sorted := make([]string, len(paths))
	copy(sorted, paths)
	sort.Strings(sorted)

	h := sha256.New()
	for _, p := range sorted {
		info, err := os.Stat(p)
		if err != nil {
			return "", fmt.Errorf("cache: stat %q: %w", p, err)
		}
		fmt.Fprintf(h, "%s:%d\n", p, info.ModTime().UnixNano())
	}
	return hex.EncodeToString(h.Sum(nil))[:16], nil
}

// Get retrieves a cached entry by key. Returns (nil, nil) on a cache miss.
func (c *Cache) Get(key string) (*Entry, error) {
	path := filepath.Join(c.dir, key+".json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache: read %q: %w", path, err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("cache: unmarshal %q: %w", path, err)
	}
	return &entry, nil
}

// Put stores a merge result under the given key.
func (c *Cache) Put(key string, result map[string]interface{}) error {
	entry := Entry{
		Key:       key,
		CreatedAt: time.Now().UTC(),
		Data:      result,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	path := filepath.Join(c.dir, key+".json")
	return os.WriteFile(path, data, 0o644)
}

// Invalidate removes the cache entry for the given key, if present.
func (c *Cache) Invalidate(key string) error {
	path := filepath.Join(c.dir, key+".json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cache: invalidate %q: %w", key, err)
	}
	return nil
}
