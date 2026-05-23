package runner_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/config"
	"github.com/logslice/logslice/internal/runner"
)

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("mustTime: %v", err)
	}
	return parsed
}

func TestRunner_Run_BasicRange(t *testing.T) {
	cfg := &config.Config{
		InputFile: "testdata/sample.log",
		Start:     mustTime(t, "2024-01-01T10:00:00Z"),
		End:       mustTime(t, "2024-01-01T10:01:00Z"),
		Format:    "raw",
	}
	var buf bytes.Buffer
	r := runner.NewWithWriter(cfg, &buf)
	count, err := r.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count == 0 {
		t.Error("expected at least one matched line")
	}
	if !strings.Contains(buf.String(), "2024-01-01T10:00") {
		t.Errorf("output missing expected timestamp prefix, got: %s", buf.String())
	}
}

func TestRunner_Run_WithPattern(t *testing.T) {
	cfg := &config.Config{
		InputFile: "testdata/sample.log",
		Pattern:   "ERROR",
		Format:    "raw",
	}
	var buf bytes.Buffer
	r := runner.NewWithWriter(cfg, &buf)
	count, err := r.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if line != "" && !strings.Contains(line, "ERROR") {
			t.Errorf("line does not match pattern: %s", line)
		}
	}
	_ = count
}

func TestRunner_Run_MissingFile(t *testing.T) {
	cfg := &config.Config{
		InputFile: "testdata/nonexistent.log",
		Format:    "raw",
	}
	var buf bytes.Buffer
	r := runner.NewWithWriter(cfg, &buf)
	_, err := r.Run()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRunner_Run_InvalidPattern(t *testing.T) {
	cfg := &config.Config{
		InputFile: "testdata/sample.log",
		Pattern:   "[invalid",
		Format:    "raw",
	}
	var buf bytes.Buffer
	r := runner.NewWithWriter(cfg, &buf)
	_, err := r.Run()
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}
