package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/confmerge/internal/loader"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoadFile_YAML(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "host: localhost\nport: 8080\n")
	out, err := loader.LoadFile(path)
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

func TestLoadFile_TOML(t *testing.T) {
	path := writeTempFile(t, "config.toml", "host = \"localhost\"\nport = 8080\n")
	out, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", out["host"])
	}
}

func TestLoadFile_UnsupportedExtension(t *testing.T) {
	path := writeTempFile(t, "config.json", `{"key": "value"}`)
	_, err := loader.LoadFile(path)
	if err == nil {
		t.Fatal("expected error for unsupported extension, got nil")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := loader.LoadFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFile_EmptyYAML(t *testing.T) {
	path := writeTempFile(t, "empty.yaml", "")
	out, err := loader.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error for empty YAML file: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map for empty YAML, got %v", out)
	}
}

func TestDetectFormat(t *testing.T) {
	cases := []struct {
		path   string
		want   loader.Format
		wantErr bool
	}{
		{"app.yaml", loader.FormatYAML, false},
		{"app.yml", loader.FormatYAML, false},
		{"app.toml", loader.FormatTOML, false},
		{"app.json", "", true},
	}
	for _, tc := range cases {
		got, err := loader.DetectFormat(tc.path)
		if tc.wantErr && err == nil {
			t.Errorf("DetectFormat(%q): expected error", tc.path)
		}
		if !tc.wantErr && got != tc.want {
			t.Errorf("DetectFormat(%q): got %q, want %q", tc.path, got, tc.want)
		}
	}
}
