package writer_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/confmerge/internal/writer"
)

var sampleData = map[string]interface{}{
	"app": map[string]interface{}{
		"name": "confmerge",
		"port": 8080,
	},
	"debug": true,
}

func TestWrite_YAML(t *testing.T) {
	var buf bytes.Buffer
	if err := writer.Write(&buf, sampleData, writer.FormatYAML); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "confmerge") {
		t.Errorf("expected YAML output to contain 'confmerge', got:\n%s", out)
	}
	if !strings.Contains(out, "port:") {
		t.Errorf("expected YAML output to contain 'port:', got:\n%s", out)
	}
}

func TestWrite_TOML(t *testing.T) {
	var buf bytes.Buffer
	if err := writer.Write(&buf, sampleData, writer.FormatTOML); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "confmerge") {
		t.Errorf("expected TOML output to contain 'confmerge', got:\n%s", out)
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := writer.Write(&buf, sampleData, writer.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"confmerge"`) {
		t.Errorf("expected JSON output to contain '\"confmerge\"', got:\n%s", out)
	}
	if !strings.Contains(out, `"debug": true`) {
		t.Errorf("expected JSON output to contain '\"debug\": true', got:\n%s", out)
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := writer.Write(&buf, sampleData, writer.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("unexpected error message: %v", err)
	}
}
