// Package auditor provides structured audit logging for confmerge pipeline
// operations. It records timestamped entries with a level (INFO, WARN, ERROR),
// a stage name identifying which pipeline phase produced the entry, and a
// human-readable message.
//
// Usage:
//
//	a := auditor.New(os.Stdout)
//	a.Info("load", "loaded 3 config files")
//	a.Warn("validate", "null value at key 'db.password'")
//	a.Error("schema", "required field 'host' missing")
//
//	for _, e := range a.Entries() {
//		fmt.Println(e.Timestamp, e.Level, e.Stage, e.Message)
//	}
package auditor
