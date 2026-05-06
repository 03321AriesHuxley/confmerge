package auditor

import (
	"fmt"
	"strings"
)

// FormatText returns a human-readable multi-line summary of all audit entries.
func FormatText(entries []Entry) string {
	if len(entries) == 0 {
		return "(no audit entries)\n"
	}
	var sb strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&sb, "[%s] %-5s %-20s %s\n",
			e.Timestamp.Format("15:04:05"), e.Level, e.Stage, e.Message)
	}
	return sb.String()
}

// Summary returns a short count-based summary string, e.g.:
// "3 entries: 1 ERROR, 1 WARN, 1 INFO"
func Summary(entries []Entry) string {
	counts := map[string]int{}
	for _, e := range entries {
		counts[e.Level]++
	}
	parts := []string{}
	for _, lvl := range []string{"ERROR", "WARN", "INFO"} {
		if n, ok := counts[lvl]; ok && n > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", n, lvl))
		}
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%d entries", len(entries))
	}
	return fmt.Sprintf("%d entries: %s", len(entries), strings.Join(parts, ", "))
}
