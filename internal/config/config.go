// Package config provides configuration types and validation for logslice.
package config

import (
	"errors"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Format represents an output format for matched log lines.
type Format string

const (
	FormatRaw  Format = "raw"
	FormatJSON Format = "json"
)

// Config holds all runtime configuration for a logslice run.
type Config struct {
	Input   string
	Output  string
	Start   time.Time
	End     time.Time
	Pattern string
	Format  Format
	Tail    bool
	PollMs  int
}

// Validate checks that the Config is well-formed and applies defaults.
func (c *Config) Validate() error {
	if c.Input == "" {
		return errors.New("input file is required")
	}
	if !c.End.IsZero() && !c.Start.IsZero() && c.End.Before(c.Start) {
		return errors.New("end time must be after start time")
	}
	if c.Format == "" {
		c.Format = FormatRaw
	}
	if c.Format != FormatRaw && c.Format != FormatJSON {
		return errors.New("format must be \"raw\" or \"json\"")
	}
	if c.Tail && c.PollMs <= 0 {
		c.PollMs = 100
	}
	if c.Pattern != "" {
		if _, err := parser.CompilePattern(c.Pattern); err != nil {
			return err
		}
	}
	return nil
}

// PollDuration returns the polling interval as a time.Duration.
func (c *Config) PollDuration() time.Duration {
	return time.Duration(c.PollMs) * time.Millisecond
}
