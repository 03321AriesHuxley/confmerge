package schema_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/confmerge/internal/schema"
)

func writeTempSchema(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "schema.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp schema: %v", err)
	}
	return p
}

func TestLoadSchema_Valid(t *testing.T) {
	p := writeTempSchema(t, "host:\n  type: string\n  required: true\nport:\n  type: int\n  required: false\n")
	s, err := schema.LoadSchema(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(s))
	}
}

func TestLoadSchema_NotFound(t *testing.T) {
	_, err := schema.LoadSchema("/nonexistent/schema.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_RequiredPresent(t *testing.T) {
	s := schema.Schema{
		"host": {Type: schema.TypeString, Required: true},
	}
	cfg := map[string]any{"host": "localhost"}
	violations := schema.Validate(cfg, s)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got: %v", violations)
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	s := schema.Schema{
		"host": {Type: schema.TypeString, Required: true},
	}
	cfg := map[string]any{}
	violations := schema.Validate(cfg, s)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(violations), violations)
	}
}

func TestValidate_WrongType(t *testing.T) {
	s := schema.Schema{
		"port": {Type: schema.TypeInt, Required: false},
	}
	cfg := map[string]any{"port": "not-an-int"}
	violations := schema.Validate(cfg, s)
	if len(violations) != 1 {
		t.Fatalf("expected 1 type violation, got %d: %v", len(violations), violations)
	}
}

func TestValidate_OptionalMissing(t *testing.T) {
	s := schema.Schema{
		"debug": {Type: schema.TypeBool, Required: false},
	}
	cfg := map[string]any{}
	violations := schema.Validate(cfg, s)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for optional missing key, got: %v", violations)
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	s := schema.Schema{
		"host": {Type: schema.TypeString, Required: true},
		"port": {Type: schema.TypeInt, Required: true},
	}
	cfg := map[string]any{}
	violations := schema.Validate(cfg, s)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d: %v", len(violations), violations)
	}
}
