package sorter_test

import (
	"testing"

	"github.com/yourorg/confmerge/internal/sorter"
)

func TestSort_AscendingKeys(t *testing.T) {
	input := map[string]any{"zebra": 1, "apple": 2, "mango": 3}
	opts := sorter.DefaultOptions()
	out := sorter.Sort(input, opts)
	keys := sorter.Keys(out, opts)
	if keys[0] != "apple" || keys[1] != "mango" || keys[2] != "zebra" {
		t.Fatalf("unexpected order: %v", keys)
	}
}

func TestSort_DescendingKeys(t *testing.T) {
	input := map[string]any{"zebra": 1, "apple": 2, "mango": 3}
	opts := sorter.Options{Recursive: false, CaseInsensitive: true, Descending: true}
	out := sorter.Sort(input, opts)
	keys := sorter.Keys(out, opts)
	if keys[0] != "zebra" || keys[1] != "mango" || keys[2] != "apple" {
		t.Fatalf("unexpected order: %v", keys)
	}
}

func TestSort_CaseInsensitive(t *testing.T) {
	input := map[string]any{"Banana": 1, "apple": 2, "Cherry": 3}
	opts := sorter.Options{Recursive: false, CaseInsensitive: true, Descending: false}
	out := sorter.Sort(input, opts)
	keys := sorter.Keys(out, opts)
	if keys[0] != "apple" || keys[1] != "Banana" || keys[2] != "Cherry" {
		t.Fatalf("unexpected case-insensitive order: %v", keys)
	}
}

func TestSort_Recursive(t *testing.T) {
	input := map[string]any{
		"z": map[string]any{"b": 1, "a": 2},
		"a": 3,
	}
	opts := sorter.Options{Recursive: true, CaseInsensitive: false, Descending: false}
	out := sorter.Sort(input, opts)
	child, ok := out["z"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map for key 'z'")
	}
	childKeys := sorter.Keys(child, opts)
	if childKeys[0] != "a" || childKeys[1] != "b" {
		t.Fatalf("nested keys not sorted: %v", childKeys)
	}
}

func TestSort_DoesNotMutateInput(t *testing.T) {
	input := map[string]any{"z": 1, "a": 2}
	opts := sorter.DefaultOptions()
	sorter.Sort(input, opts)
	if _, ok := input["z"]; !ok {
		t.Fatal("original map was mutated")
	}
}

func TestSort_NilInput(t *testing.T) {
	opts := sorter.DefaultOptions()
	out := sorter.Sort(nil, opts)
	if out != nil {
		t.Fatalf("expected nil output for nil input, got %v", out)
	}
}
