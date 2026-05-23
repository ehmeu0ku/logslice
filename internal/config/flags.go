package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// ParseFlags parses command-line arguments and returns a Config.
func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	input := fs.String("input", "", "path to log file (required)")
	output := fs.String("output", "", "path to output file (default: stdout)")
	startStr := fs.String("start", "", "start timestamp (optional)")
	endStr := fs.String("end", "", "end timestamp (optional)")
	pattern := fs.String("pattern", "", "regex pattern to filter lines (optional)")
	format := fs.String("format", "raw", "output format: raw or json")
	tailMode := fs.Bool("tail", false, "tail the file for new lines")
	pollMs := fs.Int("poll-ms", 100, "polling interval in milliseconds (tail mode)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		Input:   *input,
		Output:  *output,
		Pattern: *pattern,
		Format:  Format(*format),
		Tail:    *tailMode,
		PollMs:  *pollMs,
	}

	if *startStr != "" {
		t, err := parseTime(*startStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start time: %w", err)
		}
		cfg.Start = t
	}
	if *endStr != "" {
		t, err := parseTime(*endStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end time: %w", err)
		}
		cfg.End = t
	}

	return cfg, nil
}

func parseTime(s string) (time.Time, error) {
	for _, fmt := range parser.SupportedFormats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised time format: %q", s)
}
