package filter_test

import (
	"testing"

	"github.com/user/confmerge/internal/filter"
)

func TestFilter_NoOpts(t *testing.T) {
	input := map[string]any{"a": 1, "b": 2}
	out := filter.Filter(input, filter.Options{})
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestFilter_ExcludeTopLevel(t *testing.T) {
	input := map[string]any{"a": 1, "b": 2, "c": 3}
	out := filter.Filter(input, filter.Options{Exclude: []string{"b"}})
	if _, ok := out["b"]; ok {
		t.Fatal("expected key 'b' to be excluded")
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestFilter_IncludeTopLevel(t *testing.T) {
	input := map[string]any{"a": 1, "b": 2, "c": 3}
	out := filter.Filter(input, filter.Options{Include: []string{"a", "c"}})
	if _, ok := out["b"]; ok {
		t.Fatal("key 'b' should not be present")
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestFilter_IncludeNestedKey(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{"host": "localhost", "port": 5432},
		"app": map[string]any{"name": "myapp"},
	}
	out := filter.Filter(input, filter.Options{Include: []string{"db.host"}})
	db, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatal("expected nested 'db' map")
	}
	if _, ok := db["host"]; !ok {
		t.Fatal("expected 'db.host' to be present")
	}
	if _, ok := db["port"]; ok {
		t.Fatal("expected 'db.port' to be excluded")
	}
	if _, ok := out["app"]; ok {
		t.Fatal("expected 'app' to be excluded")
	}
}

func TestFilter_ExcludeNestedKey(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{"host": "localhost", "password": "secret"},
	}
	out := filter.Filter(input, filter.Options{Exclude: []string{"db.password"}})
	db, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatal("expected nested 'db' map")
	}
	if _, ok := db["password"]; ok {
		t.Fatal("expected 'db.password' to be excluded")
	}
	if _, ok := db["host"]; !ok {
		t.Fatal("expected 'db.host' to be present")
	}
}

func TestFilter_DoesNotMutateInput(t *testing.T) {
	input := map[string]any{"a": 1, "b": 2}
	_ = filter.Filter(input, filter.Options{Exclude: []string{"a"}})
	if _, ok := input["a"]; !ok {
		t.Fatal("Filter must not mutate the original map")
	}
}
