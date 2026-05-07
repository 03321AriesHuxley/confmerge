package flattener_test

import (
	"testing"

	"github.com/yourorg/confmerge/internal/flattener"
)

func TestFlatten_SimpleNested(t *testing.T) {
	src := map[string]any{
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
	}
	out := flattener.Flatten(src, flattener.DefaultOptions())
	if out["database.host"] != "localhost" {
		t.Errorf("expected database.host=localhost, got %v", out["database.host"])
	}
	if out["database.port"] != 5432 {
		t.Errorf("expected database.port=5432, got %v", out["database.port"])
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	src := map[string]any{"a": map[string]any{"b": "v"}}
	out := flattener.Flatten(src, flattener.Options{Separator: "__"})
	if out["a__b"] != "v" {
		t.Errorf("expected a__b=v, got %v", out["a__b"])
	}
}

func TestFlatten_MaxDepth(t *testing.T) {
	src := map[string]any{
		"a": map[string]any{
			"b": map[string]any{"c": "deep"},
		},
	}
	out := flattener.Flatten(src, flattener.Options{Separator: ".", MaxDepth: 1})
	// At depth 1 the nested map {"c":"deep"} should be kept as-is.
	if _, ok := out["a.b"]; !ok {
		t.Errorf("expected key a.b to be present")
	}
	if _, ok := out["a.b.c"]; ok {
		t.Errorf("expected a.b.c to NOT be flattened beyond MaxDepth")
	}
}

func TestFlatten_ScalarAtRoot(t *testing.T) {
	src := map[string]any{"key": "value"}
	out := flattener.Flatten(src, flattener.DefaultOptions())
	if out["key"] != "value" {
		t.Errorf("expected key=value, got %v", out["key"])
	}
}

func TestUnflatten_RoundTrip(t *testing.T) {
	orig := map[string]any{
		"server": map[string]any{
			"host": "0.0.0.0",
			"tls":  map[string]any{"enabled": true},
		},
	}
	flat := flattener.Flatten(orig, flattener.DefaultOptions())
	restored, err := flattener.Unflatten(flat, flattener.DefaultOptions())
	if err != nil {
		t.Fatalf("Unflatten error: %v", err)
	}
	srv, ok := restored["server"].(map[string]any)
	if !ok {
		t.Fatalf("expected server to be a map")
	}
	if srv["host"] != "0.0.0.0" {
		t.Errorf("expected server.host=0.0.0.0, got %v", srv["host"])
	}
	tls, ok := srv["tls"].(map[string]any)
	if !ok {
		t.Fatalf("expected server.tls to be a map")
	}
	if tls["enabled"] != true {
		t.Errorf("expected server.tls.enabled=true, got %v", tls["enabled"])
	}
}

func TestUnflatten_ConflictReturnsError(t *testing.T) {
	flat := map[string]any{
		"a.b":   "scalar",
		"a.b.c": "conflict",
	}
	_, err := flattener.Unflatten(flat, flattener.DefaultOptions())
	if err == nil {
		t.Error("expected an error for conflicting keys, got nil")
	}
}
