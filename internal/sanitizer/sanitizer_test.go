package sanitizer_test

import (
	"testing"

	"github.com/yourorg/confmerge/internal/sanitizer"
)

func TestSanitize_TrimStrings(t *testing.T) {
	input := map[string]any{"host": "  localhost  ", "port": "8080"}
	opts := sanitizer.Options{TrimStrings: true}
	out := sanitizer.Sanitize(input, opts)
	if out["host"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", out["host"])
	}
	if out["port"] != "8080" {
		t.Errorf("expected '8080', got %q", out["port"])
	}
}

func TestSanitize_DropNulls(t *testing.T) {
	input := map[string]any{"key": "value", "empty": nil}
	opts := sanitizer.Options{DropNulls: true}
	out := sanitizer.Sanitize(input, opts)
	if _, ok := out["empty"]; ok {
		t.Error("expected nil key to be dropped")
	}
	if out["key"] != "value" {
		t.Errorf("expected 'value', got %v", out["key"])
	}
}

func TestSanitize_CoerceNumbers(t *testing.T) {
	input := map[string]any{"count": "42", "ratio": "3.14", "name": "alice"}
	opts := sanitizer.Options{TrimStrings: true, CoerceNumbers: true}
	out := sanitizer.Sanitize(input, opts)
	if out["count"] != int64(42) {
		t.Errorf("expected int64(42), got %v (%T)", out["count"], out["count"])
	}
	if out["ratio"] != float64(3.14) {
		t.Errorf("expected float64(3.14), got %v (%T)", out["ratio"], out["ratio"])
	}
	if out["name"] != "alice" {
		t.Errorf("expected 'alice', got %v", out["name"])
	}
}

func TestSanitize_NestedMap(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{
			"host": "  db.local  ",
			"pass": nil,
		},
	}
	opts := sanitizer.DefaultOptions()
	out := sanitizer.Sanitize(input, opts)
	db, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map under 'db'")
	}
	if db["host"] != "db.local" {
		t.Errorf("expected 'db.local', got %q", db["host"])
	}
	if _, exists := db["pass"]; exists {
		t.Error("expected nil 'pass' key to be dropped")
	}
}

func TestSanitize_DoesNotMutateInput(t *testing.T) {
	input := map[string]any{"key": "  value  ", "nil_key": nil}
	opts := sanitizer.DefaultOptions()
	_ = sanitizer.Sanitize(input, opts)
	if input["key"] != "  value  " {
		t.Error("original input was mutated")
	}
	if _, ok := input["nil_key"]; !ok {
		t.Error("original nil key was removed from input")
	}
}

func TestSanitize_SliceDropNulls(t *testing.T) {
	input := map[string]any{"items": []any{"a", nil, "b"}}
	opts := sanitizer.Options{DropNulls: true}
	out := sanitizer.Sanitize(input, opts)
	items, ok := out["items"].([]any)
	if !ok {
		t.Fatal("expected slice under 'items'")
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items after dropping nil, got %d", len(items))
	}
}
