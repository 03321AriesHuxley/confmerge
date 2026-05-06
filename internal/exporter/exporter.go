package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format represents a supported export format.
type Format string

const (
	FormatYAML Format = "yaml"
	FormatJSON Format = "json"
	FormatEnv  Format = "env"
)

// Export serializes the merged config map into the requested format and writes
// it to w. Supported formats: yaml, json, env.
func Export(w io.Writer, data map[string]any, format Format) error {
	switch format {
	case FormatYAML:
		return exportYAML(w, data)
	case FormatJSON:
		return exportJSON(w, data)
	case FormatEnv:
		return exportEnv(w, data, "")
	default:
		return fmt.Errorf("exporter: unsupported format %q", format)
	}
}

func exportYAML(w io.Writer, data map[string]any) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("exporter: yaml encode: %w", err)
	}
	return enc.Close()
}

func exportJSON(w io.Writer, data map[string]any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("exporter: json encode: %w", err)
	}
	return nil
}

// exportEnv flattens the map into KEY=VALUE lines using uppercase keys.
// Nested keys are joined with underscores.
func exportEnv(w io.Writer, data map[string]any, prefix string) error {
	for _, k := range sortedKeys(data) {
		v := data[k]
		key := strings.ToUpper(k)
		if prefix != "" {
			key = prefix + "_" + key
		}
		switch val := v.(type) {
		case map[string]any:
			if err := exportEnv(w, val, key); err != nil {
				return err
			}
		default:
			if _, err := fmt.Fprintf(w, "%s=%v\n", key, val); err != nil {
				return fmt.Errorf("exporter: env write: %w", err)
			}
		}
	}
	return nil
}
