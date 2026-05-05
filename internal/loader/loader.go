// Package loader provides functionality for reading and parsing
// YAML and TOML configuration files into generic map structures.
package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Format represents a supported configuration file format.
type Format string

const (
	FormatYAML Format = "yaml"
	FormatTOML Format = "toml"
)

// DetectFormat infers the file format from the file extension.
// Returns an error if the extension is not recognized.
func DetectFormat(path string) (Format, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	switch ext {
	case "yaml", "yml":
		return FormatYAML, nil
	case "toml":
		return FormatTOML, nil
	default:
		return "", fmt.Errorf("unsupported file extension: %q", ext)
	}
}

// LoadFile reads a YAML or TOML file at the given path and returns
// its contents as a map[string]interface{}. The format is auto-detected
// from the file extension.
func LoadFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", path, err)
	}

	fmt, err := DetectFormat(path)
	if err != nil {
		return nil, err
	}

	return parse(data, fmt)
}

// parse decodes raw bytes into a map using the specified format.
func parse(data []byte, fmt Format) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	switch fmt {
	case FormatYAML:
		if err := yaml.Unmarshal(data, &out); err != nil {
			return nil, fmt.Errorf("parsing YAML: %w", err)
		}
	case FormatTOML:
		if _, err := toml.Decode(string(data), &out); err != nil {
			return nil, fmt.Errorf("parsing TOML: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown format: %q", fmt)
	}

	return out, nil
}
