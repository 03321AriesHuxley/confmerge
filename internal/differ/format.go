package differ

import (
	"fmt"
	"io"
	"strings"
)

// FormatText writes a human-readable diff of entries to the given writer.
// Added lines are prefixed with '+', removed with '-', modified show old->new,
// and unchanged lines are prefixed with ' '.
func FormatText(w io.Writer, entries []Entry) error {
	for _, e := range entries {
		var line string
		switch e.Type {
		case Added:
			line = fmt.Sprintf("+ %s: %v", e.Key, e.NewValue)
		case Removed:
			line = fmt.Sprintf("- %s: %v", e.Key, e.OldValue)
		case Modified:
			line = fmt.Sprintf("~ %s: %v -> %v", e.Key, e.OldValue, e.NewValue)
		case Unchanged:
			line = fmt.Sprintf("  %s: %v", e.Key, e.OldValue)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

// Summary returns a short string summarising the diff counts.
func Summary(entries []Entry) string {
	counts := map[DiffType]int{}
	for _, e := range entries {
		counts[e.Type]++
	}
	parts := []string{}
	if n := counts[Added]; n > 0 {
		parts = append(parts, fmt.Sprintf("%d added", n))
	}
	if n := counts[Removed]; n > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", n))
	}
	if n := counts[Modified]; n > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", n))
	}
	if n := counts[Unchanged]; n > 0 {
		parts = append(parts, fmt.Sprintf("%d unchanged", n))
	}
	if len(parts) == 0 {
		return "no differences"
	}
	return strings.Join(parts, ", ")
}
