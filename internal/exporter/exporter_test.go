package exporter

import (
	"bytes"
	"strings"
	"testing"
)

func TestExport_YAML(t *testing.T) {
	data := map[string]any{
		"app": map[string]any{
			"name": "confmerge",
			"port": 8080,
		},
	}
	var buf bytes.Buffer
	if err := Export(&buf, data, FormatYAML); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "name: confmerge") {
		t.Errorf("expected YAML to contain 'name: confmerge', got:\n%s", out)
	}
	if !strings.Contains(out, "port: 8080") {
		t.Errorf("expected YAML to contain 'port: 8080', got:\n%s", out)
	}
}

func TestExport_JSON(t *testing.T) {
	data := map[string]any{"key": "value", "num": 42}
	var buf bytes.Buffer
	if err := Export(&buf, data, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"key": "value"`) {
		t.Errorf("expected JSON output to contain key/value, got:\n%s", out)
	}
}

func TestExport_Env_Flat(t *testing.T) {
	data := map[string]any{"host": "localhost", "port": 5432}
	var buf bytes.Buffer
	if err := Export(&buf, data, FormatEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "HOST=localhost") {
		t.Errorf("expected HOST=localhost in env output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=5432") {
		t.Errorf("expected PORT=5432 in env output, got:\n%s", out)
	}
}

func TestExport_Env_Nested(t *testing.T) {
	data := map[string]any{
		"db": map[string]any{
			"host": "db.local",
			"port": 5432,
		},
	}
	var buf bytes.Buffer
	if err := Export(&buf, data, FormatEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=db.local") {
		t.Errorf("expected DB_HOST=db.local, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432, got:\n%s", out)
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, map[string]any{}, Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}
