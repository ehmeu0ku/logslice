package config

import (
	"testing"
)

func TestParseFlags_Minimal(t *testing.T) {
	cfg, err := ParseFlags([]string{"-f", "app.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.InputFile != "app.log" {
		t.Errorf("expected input file app.log, got %q", cfg.InputFile)
	}
	if cfg.OutputFmt != "raw" {
		t.Errorf("expected default format raw, got %q", cfg.OutputFmt)
	}
}

func TestParseFlags_WithTimeRange(t *testing.T) {
	cfg, err := ParseFlags([]string{
		"-f", "app.log",
		"-start", "2024-01-01T00:00:00Z",
		"-end", "2024-01-02T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Start == nil || cfg.End == nil {
		t.Fatal("expected start and end to be set")
	}
	if !cfg.End.After(*cfg.Start) {
		t.Error("expected end to be after start")
	}
}

func TestParseFlags_WithPattern(t *testing.T) {
	cfg, err := ParseFlags([]string{"-f", "app.log", "-pattern", "ERROR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pattern != "ERROR" {
		t.Errorf("expected pattern ERROR, got %q", cfg.Pattern)
	}
}

func TestParseFlags_InvalidStartTime(t *testing.T) {
	_, err := ParseFlags([]string{"-f", "app.log", "-start", "not-a-time"})
	if err == nil {
		t.Fatal("expected error for invalid start time")
	}
}

func TestParseFlags_JSONFormat(t *testing.T) {
	cfg, err := ParseFlags([]string{"-f", "app.log", "-format", "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFmt != FormatJSON {
		t.Errorf("expected json format, got %q", cfg.OutputFmt)
	}
}

func TestParseFlags_Summary(t *testing.T) {
	cfg, err := ParseFlags([]string{"-f", "app.log", "-summary"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.ShowSummary {
		t.Error("expected ShowSummary to be true")
	}
}
