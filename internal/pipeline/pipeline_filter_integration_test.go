package pipeline_test

import (
	"testing"

	"github.com/user/confmerge/internal/filter"
)

// TestFilter_PipelineKeyExclusion verifies that the filter stage correctly
// removes sensitive keys before the result reaches downstream pipeline steps.
func TestFilter_PipelineKeyExclusion(t *testing.T) {
	input := map[string]any{
		"app": map[string]any{
			"name":   "confmerge",
			"secret": "topsecret",
		},
		"database": map[string]any{
			"host":     "localhost",
			"password": "hunter2",
		},
	}

	result := filter.Filter(input, filter.Options{
		Exclude: []string{"app.secret", "database.password"},
	})

	app, ok := result["app"].(map[string]any)
	if !ok {
		t.Fatal("expected app map")
	}
	if _, found := app["secret"]; found {
		t.Error("app.secret should have been excluded")
	}

	db, ok := result["database"].(map[string]any)
	if !ok {
		t.Fatal("expected database map")
	}
	if _, found := db["password"]; found {
		t.Error("database.password should have been excluded")
	}
	if db["host"] != "localhost" {
		t.Errorf("database.host should be preserved, got %v", db["host"])
	}
}

// TestFilter_PipelineIncludeOnly verifies that only whitelisted keys survive.
func TestFilter_PipelineIncludeOnly(t *testing.T) {
	input := map[string]any{
		"a": 1,
		"b": 2,
		"c": map[string]any{"x": 10, "y": 20},
	}

	result := filter.Filter(input, filter.Options{
		Include: []string{"a", "c.x"},
	})

	if result["a"] != 1 {
		t.Errorf("expected a=1, got %v", result["a"])
	}
	if _, ok := result["b"]; ok {
		t.Error("b should not be included")
	}
	c, ok := result["c"].(map[string]any)
	if !ok {
		t.Fatal("expected c map")
	}
	if c["x"] != 10 {
		t.Errorf("expected c.x=10, got %v", c["x"])
	}
	if _, ok := c["y"]; ok {
		t.Error("c.y should not be included")
	}
}
