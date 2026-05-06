// Package profiler provides lightweight pipeline profiling for confmerge.
//
// It allows each stage of the merge pipeline (resolve, load, merge, transform,
// write) to be timed independently using a deferred tracking pattern:
//
//	p := profiler.New()
//	defer p.Track("load")()
//
// After all stages complete, Print writes a formatted timing summary to any
// io.Writer, and Stages returns the raw slice for programmatic inspection.
package profiler
