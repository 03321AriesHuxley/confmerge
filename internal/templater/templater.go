// Package templater provides Go template rendering support for merged config values.
// It allows values in config maps to reference other keys using Go template syntax,
// enabling self-referential configs such as:
//
//	base_url: "https://example.com"
//	api_url:  "{{ .base_url }}/api/v1"
package templater

import (
	"bytes"
	"fmt"
	"text/template"
)

// Render walks the config map and evaluates any string values that contain
// Go template expressions. The entire (flat) map is available as template data,
// so only top-level keys can be referenced. Nested maps are passed through
// unchanged. Returns a new map; the original is not mutated.
func Render(data map[string]interface{}) (map[string]interface{}, error) {
	// Build a flat string-only view for template data.
	tplData := make(map[string]interface{}, len(data))
	for k, v := range data {
		tplData[k] = v
	}

	out := make(map[string]interface{}, len(data))
	for k, v := range data {
		switch val := v.(type) {
		case string:
			rendered, err := renderString(val, tplData)
			if err != nil {
				return nil, fmt.Errorf("templater: key %q: %w", k, err)
			}
			out[k] = rendered
		case map[string]interface{}:
			nested, err := Render(val)
			if err != nil {
				return nil, fmt.Errorf("templater: key %q: %w", k, err)
			}
			out[k] = nested
		default:
			out[k] = v
		}
	}
	return out, nil
}

func renderString(s string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("").Option("missingkey=error").Parse(s)
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute: %w", err)
	}
	return buf.String(), nil
}
