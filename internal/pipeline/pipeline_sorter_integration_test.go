package pipeline_test

import (
	"testing"

	"github.com/yourorg/confmerge/internal/sorter"
)

// TestSorter_PipelineStageOrder verifies that applying the sorter after a
// merge produces stable key ordering, which is the expected behaviour when
// the sorter is wired into the pipeline.
func TestSorter_PipelineStageOrder(t *testing.T) {
	merged := map[string]any{
		"zoo":    "animal",
		"alpha":  42,
		"middle": map[string]any{"z": true, "a": false},
	}

	opts := sorter.DefaultOptions()
	sorted := sorter.Sort(merged, opts)
	keys := sorter.Keys(sorted, opts)

	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "alpha" || keys[1] != "middle" || keys[2] != "zoo" {
		t.Fatalf("top-level keys not sorted correctly: %v", keys)
	}

	child, ok := sorted["middle"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map")
	}
	childKeys := sorter.Keys(child, opts)
	if childKeys[0] != "a" || childKeys[1] != "z" {
		t.Fatalf("nested keys not sorted: %v", childKeys)
	}
}
