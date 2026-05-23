package config

import (
	"errors"
	"time"
)

// Format represents the output format for log lines.
type Format string

const (
	FormatRaw   Format = "raw"
	FormatJSON  Format = "json"
	FormatCount Format = "count"
)

// Config holds all runtime configuration for a logslice run.
type Config struct {
	InputFile   string
	OutputFile  string
	Start       *time.Time
	End         *time.Time
	Pattern     string
	OutputFmt   Format
	TimeFmt     string
	ShowSummary bool
}

// Validate checks that the configuration is valid and consistent.
func (c *Config) Validate() error {
	if c.InputFile == "" {
		return errors.New("input file is required")
	}
	if c.Start != nil && c.End != nil && c.End.Before(*c.Start) {
		return errors.New("end time must be after start time")
	}
	switch c.OutputFmt {
	case FormatRaw, FormatJSON, FormatCount:
		// valid
	case "":
		c.OutputFmt = FormatRaw
	default:
		return errors.New("unsupported output format: " + string(c.OutputFmt))
	}
	return nil
}

// HasTimeRange returns true if at least one time bound is set.
func (c *Config) HasTimeRange() bool {
	return c.Start != nil || c.End != nil
}
