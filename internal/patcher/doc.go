// Package patcher provides patch-operation support for config maps.
//
// A Patch describes a single mutation using an Op (set, delete, or merge)
// and a dot-notation Path that addresses a key within a nested
// map[string]interface{} config structure.
//
// Example usage:
//
//	patches := []patcher.Patch{
//		{Op: patcher.OpSet,    Path: "server.port",  Value: 9090},
//		{Op: patcher.OpDelete, Path: "debug"},
//		{Op: patcher.OpMerge,  Path: "app",          Value: map[string]interface{}{"env": "prod"}},
//	}
//	cfg, err := patcher.Apply(cfg, patches)
//
// Patches are applied in order; later patches take precedence over earlier ones.
package patcher
