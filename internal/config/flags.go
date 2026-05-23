package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// ParseFlags parses command-line flags into a Config.
func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	input := fs.String("f", "", "input log file (required)")
	output := fs.String("o", "", "output file (default: stdout)")
	startStr := fs.String("start", "", "start timestamp filter")
	endStr := fs.String("end", "", "end timestamp filter")
	pattern := fs.String("pattern", "", "regex pattern to filter log lines")
	format := fs.String("format", "raw", "output format: raw, json, count")
	timeFmt := fs.String("timefmt", "", "custom time format (Go layout)")
	summary := fs.Bool("summary", false, "print summary after processing")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		InputFile:   *input,
		OutputFile:  *output,
		Pattern:     *pattern,
		OutputFmt:   Format(*format),
		TimeFmt:     *timeFmt,
		ShowSummary: *summary,
	}

	if *startStr != "" {
		t, err := parseTime(*startStr, *timeFmt)
		if err != nil {
			return nil, fmt.Errorf("invalid start time: %w", err)
		}
		cfg.Start = &t
	}

	if *endStr != "" {
		t, err := parseTime(*endStr, *timeFmt)
		if err != nil {
			return nil, fmt.Errorf("invalid end time: %w", err)
		}
		cfg.End = &t
	}

	return cfg, nil
}

func parseTime(s, fmt string) (time.Time, error) {
	if fmt != "" {
		return parser.ParseTimestampWithFormat(s, fmt)
	}
	return parser.ParseTimestamp(s)
}
