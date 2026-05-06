package auditor

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time
	Level     string
	Stage     string
	Message   string
}

// Auditor records structured audit entries for pipeline operations.
type Auditor struct {
	w       io.Writer
	entries []Entry
}

// New creates a new Auditor writing to the given writer.
// If w is nil, os.Stderr is used.
func New(w io.Writer) *Auditor {
	if w == nil {
		w = os.Stderr
	}
	return &Auditor{w: w}
}

// Log records an audit entry with the given level, stage, and message.
func (a *Auditor) Log(level, stage, message string) {
	e := Entry{
		Timestamp: time.Now(),
		Level:     level,
		Stage:     stage,
		Message:   message,
	}
	a.entries = append(a.entries, e)
	fmt.Fprintf(a.w, "[%s] %-5s %-20s %s\n",
		e.Timestamp.Format(time.RFC3339), e.Level, e.Stage, e.Message)
}

// Info logs an informational audit entry.
func (a *Auditor) Info(stage, message string) {
	a.Log("INFO", stage, message)
}

// Warn logs a warning audit entry.
func (a *Auditor) Warn(stage, message string) {
	a.Log("WARN", stage, message)
}

// Error logs an error audit entry.
func (a *Auditor) Error(stage, message string) {
	a.Log("ERROR", stage, message)
}

// Entries returns a copy of all recorded audit entries.
func (a *Auditor) Entries() []Entry {
	copy := make([]Entry, len(a.entries))
	for i, e := range a.entries {
		copy[i] = e
	}
	return copy
}

// Count returns the total number of recorded entries.
func (a *Auditor) Count() int {
	return len(a.entries)
}
