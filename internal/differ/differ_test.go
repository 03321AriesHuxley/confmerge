package differ

import (
	"testing"
)

func TestDiff_AddedKey(t *testing.T) {
	base := map[string]interface{}{"a": 1}
	override := map[string]interface{}{"a": 1, "b": 2}
	entries := Diff(base, override)
	if !containsEntry(entries, "b", Added) {
		t.Errorf("expected 'b' to be Added")
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	base := map[string]interface{}{"a": 1, "b": 2}
	override := map[string]interface{}{"a": 1}
	entries := Diff(base, override)
	if !containsEntry(entries, "b", Removed) {
		t.Errorf("expected 'b' to be Removed")
	}
}

func TestDiff_ModifiedKey(t *testing.T) {
	base := map[string]interface{}{"a": "old"}
	override := map[string]interface{}{"a": "new"}
	entries := Diff(base, override)
	if !containsEntry(entries, "a", Modified) {
		t.Errorf("expected 'a' to be Modified")
	}
}

func TestDiff_UnchangedKey(t *testing.T) {
	base := map[string]interface{}{"x": 42}
	override := map[string]interface{}{"x": 42}
	entries := Diff(base, override)
	if !containsEntry(entries, "x", Unchanged) {
		t.Errorf("expected 'x' to be Unchanged")
	}
}

func TestDiff_NestedModified(t *testing.T) {
	base := map[string]interface{}{
		"db": map[string]interface{}{"host": "localhost", "port": 5432},
	}
	override := map[string]interface{}{
		"db": map[string]interface{}{"host": "remotehost", "port": 5432},
	}
	entries := Diff(base, override)
	if !containsEntry(entries, "db.host", Modified) {
		t.Errorf("expected 'db.host' to be Modified")
	}
	if !containsEntry(entries, "db.port", Unchanged) {
		t.Errorf("expected 'db.port' to be Unchanged")
	}
}

func TestDiff_EmptyMaps(t *testing.T) {
	entries := Diff(map[string]interface{}{}, map[string]interface{}{})
	if len(entries) != 0 {
		t.Errorf("expected no entries for empty maps, got %d", len(entries))
	}
}

func containsEntry(entries []Entry, key string, dt DiffType) bool {
	for _, e := range entries {
		if e.Key == key && e.Type == dt {
			return true
		}
	}
	return false
}
