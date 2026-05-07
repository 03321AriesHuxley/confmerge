// Package snapshotter provides point-in-time capture and retrieval of
// configuration maps produced by the confmerge pipeline.
//
// Snapshots are persisted as JSON files in a configurable directory, allowing
// operators to roll back to a previous configuration state or diff two
// historical configurations using the differ package.
//
// Basic usage:
//
//	s, err := snapshotter.New("/var/lib/confmerge/snapshots")
//	if err != nil { ... }
//
//	// Save the current merged config.
//	if err := s.Save("v1.2.0", mergedMap); err != nil { ... }
//
//	// Reload it later.
//	snap, err := s.Load("v1.2.0")
package snapshotter
