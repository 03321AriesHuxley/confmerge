package encoder_test

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/your-org/confmerge/internal/encoder"
)

var sample = map[string]any{
	"app": map[string]any{
		"name": "confmerge",
		"port": 8080,
	},
	"debug": true,
}

func TestEncode_YAML(t *testing.T) {
	out, err := encoder.Encode(sample, encoder.FormatYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded map[string]any
	if err := yaml.Unmarshal(out, &decoded); err != nil {
		t.Fatalf("output is not valid YAML: %v", err)
	}
	if decoded["debug"] != true {
		t.Errorf("expected debug=true, got %v", decoded["debug"])
	}
}

func TestEncode_JSON(t *testing.T) {
	out, err := encoder.Encode(sample, encoder.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(out, &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	app, ok := decoded["app"].(map[string]any)
	if !ok {
		t.Fatal("expected app to be a map")
	}
	if app["name"] != "confmerge" {
		t.Errorf("expected app.name=confmerge, got %v", app["name"])
	}
}

func TestEncode_TOML(t *testing.T) {
	out, err := encoder.Encode(sample, encoder.FormatTOML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "[app]") {
		t.Errorf("expected TOML output to contain [app] section, got:\n%s", out)
	}
}

func TestEncode_UnsupportedFormat(t *testing.T) {
	_, err := encoder.Encode(sample, encoder.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEncode_EmptyMap(t *testing.T) {
	out, err := encoder.Encode(map[string]any{}, encoder.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(string(out)) != "{}" {
		t.Errorf("expected empty JSON object, got: %s", out)
	}
}
