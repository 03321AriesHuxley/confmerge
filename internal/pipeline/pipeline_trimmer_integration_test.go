package pipeline_test

import (
	"testing"

	"github.com/yourorg/confmerge/internal/trimmer"
)

// TestTrimmer_PipelineWhitespaceCleanup verifies that the trimmer correctly
// cleans up keys and values that might arrive from raw YAML/TOML parsing
// before further pipeline stages process them.
func TestTrimmer_PipelineWhitespaceCleanup(t *testing.T) {
	raw := map[string]any{
		" database ": map[string]any{
			" host ": "  127.0.0.1  ",
			" port ": 5432,
		},
		"app": map[string]any{
			"name": "  confmerge  ",
		},
	}

	out := trimmer.Trim(raw, trimmer.DefaultOptions())

	db, ok := out["database"].(map[string]any)
	if !ok {
		t.Fatal("expected 'database' key after trim")
	}
	if db["host"] != "127.0.0.1" {
		t.Errorf("host: expected '127.0.0.1', got %q", db["host"])
	}
	if db["port"] != 5432 {
		t.Errorf("port: expected 5432, got %v", db["port"])
	}

	app, ok := out["app"].(map[string]any)
	if !ok {
		t.Fatal("expected 'app' key after trim")
	}
	if app["name"] != "confmerge" {
		t.Errorf("name: expected 'confmerge', got %q", app["name"])
	}
}

func TestTrimmer_PipelinePrefixStrip(t *testing.T) {
	raw := map[string]any{
		"endpoints": []any{
			"https://api.example.com",
			"https://cdn.example.com",
		},
	}

	out := trimmer.Trim(raw, trimmer.Options{
		TrimStringValues: false,
		TrimPrefix:       "https://",
	})

	eps, ok := out["endpoints"].([]any)
	if !ok || len(eps) != 2 {
		t.Fatal("expected slice of length 2")
	}
	if eps[0] != "api.example.com" {
		t.Errorf("expected 'api.example.com', got %q", eps[0])
	}
	if eps[1] != "cdn.example.com" {
		t.Errorf("expected 'cdn.example.com', got %q", eps[1])
	}
}
