package templater_test

import (
	"testing"

	"github.com/yourorg/confmerge/internal/templater"
)

func TestRender_NoTemplates(t *testing.T) {
	input := map[string]interface{}{
		"host": "localhost",
		"port": 8080,
	}
	out, err := templater.Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", out["host"])
	}
	if out["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", out["port"])
	}
}

func TestRender_SimpleReference(t *testing.T) {
	input := map[string]interface{}{
		"base_url": "https://example.com",
		"api_url":  `{{ index . "base_url" }}/api/v1`,
	}
	out, err := templater.Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://example.com/api/v1"
	if out["api_url"] != want {
		t.Errorf("expected %q, got %q", want, out["api_url"])
	}
}

func TestRender_NestedMap(t *testing.T) {
	input := map[string]interface{}{
		"scheme": "https",
		"db": map[string]interface{}{
			"dsn": `{{ index . "scheme" }}://db.example.com`,
		},
	}
	out, err := templater.Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db, ok := out["db"].(map[string]interface{})
	if !ok {
		t.Fatal("expected nested map for 'db'")
	}
	want := "https://db.example.com"
	if db["dsn"] != want {
		t.Errorf("expected %q, got %q", want, db["dsn"])
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	input := map[string]interface{}{
		"bad": "{{ .unclosed",
	}
	_, err := templater.Render(input)
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
}

func TestRender_MissingKey(t *testing.T) {
	input := map[string]interface{}{
		"url": `{{ index . "nonexistent" }}/path`,
	}
	_, err := templater.Render(input)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRender_DoesNotMutateInput(t *testing.T) {
	input := map[string]interface{}{
		"base": "http://example.com",
		"full": `{{ index . "base" }}/v1`,
	}
	_, err := templater.Render(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input["full"] != `{{ index . "base" }}/v1` {
		t.Error("Render mutated the input map")
	}
}
