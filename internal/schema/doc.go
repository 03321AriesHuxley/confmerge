// Package schema provides schema definition loading and validation for
// merged configuration maps.
//
// A schema is a YAML file that maps top-level config keys to field
// definitions, each specifying an expected type and whether the field
// is required. Example schema file:
//
//	host:
//	  type: string
//	  required: true
//	port:
//	  type: int
//	  required: false
//
// Use LoadSchema to parse such a file, then Validate to check a
// config map produced by the merger against it.
package schema
