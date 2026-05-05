// Package loader handles reading configuration files from disk and
// deserializing them into a uniform map[string]interface{} representation
// suitable for deep-merging by the merger package.
//
// # Supported Formats
//
// The loader currently supports:
//   - YAML (.yaml, .yml)
//   - TOML (.toml)
//
// Format detection is automatic and based on the file extension.
//
// # Usage
//
//	out, err := loader.LoadFile("config.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The returned map can be passed directly to merger.Merge for layered
// override processing.
package loader
