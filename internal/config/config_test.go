package config

import (
	"testing"
	"time"
)

func TestConfig_Validate_MissingInput(t *testing.T) {
	c := &Config{}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing input file")
	}
}

func TestConfig_Validate_InvalidTimeRange(t *testing.T) {
	start := time.Now()
	end := start.Add(-time.Hour)
	c := &Config{
		InputFile: "test.log",
		Start:     &start,
		End:       &end,
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for end before start")
	}
}

func TestConfig_Validate_InvalidFormat(t *testing.T) {
	c := &Config{
		InputFile: "test.log",
		OutputFmt: Format("xml"),
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestConfig_Validate_DefaultsFormat(t *testing.T) {
	c := &Config{InputFile: "test.log"}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.OutputFmt != FormatRaw {
		t.Errorf("expected default format %q, got %q", FormatRaw, c.OutputFmt)
	}
}

func TestConfig_Validate_ValidFull(t *testing.T) {
	start := time.Now().Add(-time.Hour)
	end := time.Now()
	c := &Config{
		InputFile:  "app.log",
		OutputFile: "out.log",
		Start:      &start,
		End:        &end,
		Pattern:    "ERROR",
		OutputFmt:  FormatJSON,
		ShowSummary: true,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfig_HasTimeRange(t *testing.T) {
	c := &Config{InputFile: "test.log"}
	if c.HasTimeRange() {
		t.Error("expected no time range")
	}
	now := time.Now()
	c.Start = &now
	if !c.HasTimeRange() {
		t.Error("expected time range with start set")
	}
}
