// Package profiler provides functionality for profiling and benchmarking
// the merge pipeline, measuring time spent in each stage.
package profiler

import (
	"fmt"
	"io"
	"time"
)

// Stage represents a named stage in the merge pipeline.
type Stage struct {
	Name     string
	Duration time.Duration
}

// Profile holds timing information for all pipeline stages.
type Profile struct {
	stages []Stage
	start  time.Time
}

// New creates a new Profile and records the overall start time.
func New() *Profile {
	return &Profile{start: time.Now()}
}

// Track records the duration of a named stage using a deferred call pattern.
// Usage: defer p.Track("load")()
func (p *Profile) Track(name string) func() {
	t := time.Now()
	return func() {
		p.stages = append(p.stages, Stage{
			Name:     name,
			Duration: time.Since(t),
		})
	}
}

// Total returns the elapsed time since the Profile was created.
func (p *Profile) Total() time.Duration {
	return time.Since(p.start)
}

// Stages returns a copy of the recorded stages.
func (p *Profile) Stages() []Stage {
	out := make([]Stage, len(p.stages))
	copy(out, p.stages)
	return out
}

// Print writes a human-readable timing report to w.
func (p *Profile) Print(w io.Writer) {
	fmt.Fprintln(w, "Pipeline profile:")
	for _, s := range p.stages {
		fmt.Fprintf(w, "  %-20s %s\n", s.Name, s.Duration.Round(time.Microsecond))
	}
	fmt.Fprintf(w, "  %-20s %s\n", "total", p.Total().Round(time.Microsecond))
}
