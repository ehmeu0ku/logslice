package parser

import (
	"regexp"
	"strings"
	"time"
)

// timestampPattern matches common timestamp prefixes at the start of a log line.
var timestampPattern = regexp.MustCompile(
	`^(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:[.,]\d+)?(?:Z|[+-]\d{2}:?\d{2})?|` +
		`\d{2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4}|` +
		`\w{3}\s+\d{2} \d{2}:\d{2}:\d{2})`,
)

// LogLine represents a single parsed log line.
type LogLine struct {
	Raw       string
	Timestamp time.Time
	HasTime   bool
}

// ParseLine extracts the timestamp (if any) from a raw log line.
func ParseLine(raw string) LogLine {
	line := LogLine{Raw: raw}

	match := timestampPattern.FindString(strings.TrimSpace(raw))
	if match == "" {
		return line
	}

	t, _, err := ParseTimestamp(match)
	if err != nil {
		return line
	}

	line.Timestamp = t
	line.HasTime = true
	return line
}

// InRange reports whether the line's timestamp falls within [start, end].
// Lines without a timestamp are never in range.
func (l LogLine) InRange(start, end time.Time) bool {
	if !l.HasTime {
		return false
	}
	return !l.Timestamp.Before(start) && !l.Timestamp.After(end)
}
