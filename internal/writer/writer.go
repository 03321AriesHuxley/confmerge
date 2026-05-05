package writer

import (
	"encoding/json"
	"fmt"	
	"io"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Format represents the supported output formats.
type Format string

const (
	FormatYAML Format = "yaml"
	FormatTOML Format = "toml"
	FormatJSON Format = "json"
)

// Write serializes the given data map to the specified format and writes it to w.
func Write(w io.Writer, data map[string]interface{}, format Format) error {
	switch strings.ToLower(string(format)) {
	case string(FormatYAML):
		return writeYAML(w, data)
	case string(FormatTOML):
		return writeTOML(w, data)
	case string(FormatJSON):
		return writeJSON(w, data)
	default:
		return fmt.Errorf("unsupported output format: %q", format)
	}
}

func writeYAML(w io.Writer, data map[string]interface{}) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("yaml encode: %w", err)
	}
	return enc.Close()
}

func writeTOML(w io.Writer, data map[string]interface{}) error {
	enc := toml.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("toml encode: %w", err)
	}
	return nil
}

func writeJSON(w io.Writer, data map[string]interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}
