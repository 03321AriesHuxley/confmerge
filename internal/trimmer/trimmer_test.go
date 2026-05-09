package trimmer_test

import (
	"testing"

	"github.com/yourorg/confmerge/internal/trimmer"
)

func TestTrim_TrimKeys(t *testing.T) {
	src := map[string]any{" key ": "value"}
	out := trimmer.Trim(src, trimmer.Options{TrimKeys: true})
	if _, ok := out["key"]; !ok {
		t.Errorf("expected key 'key' after trimming, got %v", out)
	}
}

func TestTrim_TrimStringValues(t *testing.T) {
	src := map[string]any{"greeting": "  hello  "}
	out := trimmer.Trim(src, trimmer.Options{TrimStringValues: true})
	if out["greeting"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["greeting"])
	}
}

func TestTrim_TrimPrefix(t *testing.T) {
	src := map[string]any{"url": "https://example.com"}
	out := trimmer.Trim(src, trimmer.Options{TrimPrefix: "https://"})
	if out["url"] != "example.com" {
		t.Errorf("expected 'example.com', got %q", out["url"])
	}
}

func TestTrim_TrimSuffix(t *testing.T) {
	src := map[string]any{"path": "/api/v1/"}
	out := trimmer.Trim(src, trimmer.Options{TrimSuffix: "/"})
	if out["path"] != "/api/v1" {
		t.Errorf("expected '/api/v1', got %q", out["path"])
	}
}

func TestTrim_NestedMap(t *testing.T) {
	src := map[string]any{
		"db": map[string]any{
			" host ": "  localhost  ",
		},
	}
	out := trimmer.Trim(src, trimmer.DefaultOptions())
	db, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatal("expected nested map under 'db'")
	}
	if db["host"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", db["host"])
	}
}

func TestTrim_SliceValues(t *testing.T) {
	src := map[string]any{"tags": []any{" alpha ", " beta "}}
	out := trimmer.Trim(src, trimmer.Options{TrimStringValues: true})
	tags, ok := out["tags"].([]any)
	if !ok || len(tags) != 2 {
		t.Fatal("expected slice of length 2")
	}
	if tags[0] != "alpha" || tags[1] != "beta" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestTrim_DoesNotMutateInput(t *testing.T) {
	src := map[string]any{" key ": "  value  "}
	_ = trimmer.Trim(src, trimmer.DefaultOptions())
	if _, ok := src[" key "]; !ok {
		t.Error("original map was mutated")
	}
	if src[" key "] != "  value  " {
		t.Error("original value was mutated")
	}
}

func TestTrim_NonStringValuesUnchanged(t *testing.T) {
	src := map[string]any{"count": 42, "enabled": true}
	out := trimmer.Trim(src, trimmer.DefaultOptions())
	if out["count"] != 42 {
		t.Errorf("expected 42, got %v", out["count"])
	}
	if out["enabled"] != true {
		t.Errorf("expected true, got %v", out["enabled"])
	}
}
