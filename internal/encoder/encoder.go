// Package encoder provides utilities for encoding config maps into
// various byte representations (YAML, TOML, JSON) without writing to disk.
package encoder

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Format represents a supported encoding format.
type Format string

const (
	FormatYAML Format = "yaml"
	FormatTOML Format = "toml"
	FormatJSON Format = "json"
)

// Encode serialises data into the requested format and returns the raw bytes.
// data must be a map[string]any or any value accepted by the underlying
// marshaller. An error is returned for unsupported formats or marshal failures.
func Encode(data any, format Format) ([]byte, error) {
	switch format {
	case FormatYAML:
		return encodeYAML(data)
	case FormatTOML:
		return encodeTOML(data)
	case FormatJSON:
		return encodeJSON(data)
	default:
		return nil, fmt.Errorf("encoder: unsupported format %q", format)
	}
}

func encodeYAML(data any) ([]byte, error) {
	out, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("encoder: yaml marshal: %w", err)
	}
	return out, nil
}

func encodeTOML(data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(data); err != nil {
		return nil, fmt.Errorf("encoder: toml marshal: %w", err)
	}
	return buf.Bytes(), nil
}

func encodeJSON(data any) ([]byte, error) {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encoder: json marshal: %w", err)
	}
	return append(out, '\n'), nil
}
