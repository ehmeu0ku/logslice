package parser

import (
	"fmt"
	"time"
)

// Common log timestamp formats to attempt parsing
var knownFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05.000000",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// ParseTimestamp attempts to parse a timestamp string using known formats.
// It returns the parsed time and the format that matched, or an error.
func ParseTimestamp(raw string) (time.Time, string, error) {
	for _, format := range knownFormats {
		t, err := time.Parse(format, raw)
		if err == nil {
			return t, format, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("unable to parse timestamp: %q", raw)
}

// ParseTimestampWithFormat parses a timestamp using a specific format.
func ParseTimestampWithFormat(raw, format string) (time.Time, error) {
	t, err := time.Parse(format, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse error with format %q: %w", format, err)
	}
	return t, nil
}

// SupportedFormats returns the list of formats logslice recognises by default.
func SupportedFormats() []string {
	result := make([]string, len(knownFormats))
	copy(result, knownFormats)
	return result
}
