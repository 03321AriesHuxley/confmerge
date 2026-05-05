// Package writer provides serialization utilities for writing merged
// configuration data to an output stream.
//
// It supports the following output formats:
//   - YAML  (FormatYAML)
//   - TOML  (FormatTOML)
//   - JSON  (FormatJSON)
//
// Usage:
//
//	var buf bytes.Buffer
//	err := writer.Write(&buf, mergedData, writer.FormatYAML)
//
// The Write function accepts any io.Writer, making it easy to write
// results to a file, stdout, or an in-memory buffer.
package writer
