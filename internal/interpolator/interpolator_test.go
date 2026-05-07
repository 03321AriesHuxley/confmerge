package interpolator_test

import (
	"testing"

	"github.com/your-org/confmerge/internal/interpolator"
)

func TestInterpolate_NoRefs(t *testing.T) {
	input := map[string]any{"host": "localhost", "port": 8080}
	out, err := interpolator.Interpolate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["host"] != "localhost" {
		t.Errorf("expected localhost, got %v", out["host"])
	}
}

func TestInterpolate_SimpleRef(t *testing.T) {
	input := map[string]any{
		"base": "http://example.com",
		"url":  "${base}/api",
	}
	out, err := interpolator.Interpolate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["url"]; got != "http://example.com/api" {
		t.Errorf("expected http://example.com/api, got %v", got)
	}
}

func TestInterpolate_NestedRef(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{"host": "db.local"},
		"dsn": "postgres://${db.host}/mydb",
	}
	out, err := interpolator.Interpolate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["dsn"]; got != "postgres://db.local/mydb" {
		t.Errorf("expected postgres://db.local/mydb, got %v", got)
	}
}

func TestInterpolate_RefInNestedMap(t *testing.T) {
	input := map[string]any{
		"env":  "production",
		"app":  map[string]any{"name": "myapp-${env}"},
	}
	out, err := interpolator.Interpolate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	app := out["app"].(map[string]any)
	if got := app["name"]; got != "myapp-production" {
		t.Errorf("expected myapp-production, got %v", got)
	}
}

func TestInterpolate_UnresolvedRef(t *testing.T) {
	input := map[string]any{"url": "${missing.key}/path"}
	_, err := interpolator.Interpolate(input)
	if err == nil {
		t.Fatal("expected error for unresolved reference, got nil")
	}
}

func TestInterpolate_DoesNotMutateInput(t *testing.T) {
	input := map[string]any{
		"base": "http://example.com",
		"url":  "${base}/api",
	}
	_, err := interpolator.Interpolate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input["url"] != "${base}/api" {
		t.Errorf("input was mutated: got %v", input["url"])
	}
}

func TestInterpolate_MultipleRefsInOneString(t *testing.T) {
	input := map[string]any{
		"proto": "https",
		"host":  "example.com",
		"url":   "${proto}://${host}/api",
	}
	out, err := interpolator.Interpolate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["url"]; got != "https://example.com/api" {
		t.Errorf("expected https://example.com/api, got %v", got)
	}
}
