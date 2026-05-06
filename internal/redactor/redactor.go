// Package redactor provides utilities for masking sensitive values
// in configuration maps before output or logging.
package redactor

import "strings"

// DefaultPatterns is the list of key substrings treated as sensitive by default.
var DefaultPatterns = []string{
	"password", "passwd", "secret", "token",
	"apikey", "api_key", "private_key", "credential",
}

// Options controls redaction behaviour.
type Options struct {
	// Patterns is the list of case-insensitive key substrings to redact.
	// If nil, DefaultPatterns is used.
	Patterns []string
	// Mask is the replacement string. Defaults to "***".
	Mask string
}

func (o Options) patterns() []string {
	if len(o.Patterns) > 0 {
		return o.Patterns
	}
	return DefaultPatterns
}

func (o Options) mask() string {
	if o.Mask != "" {
		return o.Mask
	}
	return "***"
}

// Redact returns a deep copy of data with sensitive values replaced by the mask.
// Keys are matched case-insensitively against the configured patterns.
func Redact(data map[string]any, opts Options) map[string]any {
	return redactMap(data, opts.patterns(), opts.mask())
}

func redactMap(m map[string]any, patterns []string, mask string) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if isSensitive(k, patterns) {
			out[k] = mask
			continue
		}
		switch val := v.(type) {
		case map[string]any:
			out[k] = redactMap(val, patterns, mask)
		default:
			out[k] = v
		}
	}
	return out
}

func isSensitive(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}
