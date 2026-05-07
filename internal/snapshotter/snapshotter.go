// Package snapshotter provides functionality for capturing and comparing
// point-in-time snapshots of merged configuration maps.
package snapshotter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a named point-in-time capture of a config map.
type Snapshot struct {
	Name      string                 `json:"name"`
	CreatedAt time.Time              `json:"created_at"`
	Data      map[string]interface{} `json:"data"`
}

// Snapshotter manages saving and loading config snapshots to disk.
type Snapshotter struct {
	dir string
}

// New creates a Snapshotter that stores snapshots in dir.
// The directory is created if it does not exist.
func New(dir string) (*Snapshotter, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshotter: create dir: %w", err)
	}
	return &Snapshotter{dir: dir}, nil
}

// Save writes a snapshot of data under the given name.
// Existing snapshots with the same name are overwritten.
func (s *Snapshotter) Save(name string, data map[string]interface{}) error {
	snap := Snapshot{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		Data:      data,
	}
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshotter: marshal: %w", err)
	}
	path := filepath.Join(s.dir, name+".json")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("snapshotter: write %s: %w", path, err)
	}
	return nil
}

// Load reads a previously saved snapshot by name.
// Returns an error if the snapshot does not exist.
func (s *Snapshotter) Load(name string) (*Snapshot, error) {
	path := filepath.Join(s.dir, name+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshotter: read %s: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return nil, fmt.Errorf("snapshotter: unmarshal: %w", err)
	}
	return &snap, nil
}

// List returns the names of all stored snapshots.
func (s *Snapshotter) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("snapshotter: list dir: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

// Delete removes a snapshot by name. It is a no-op if the snapshot does not exist.
func (s *Snapshotter) Delete(name string) error {
	path := filepath.Join(s.dir, name+".json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("snapshotter: delete %s: %w", path, err)
	}
	return nil
}
