package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/output"
)

func TestWriter_RawFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatRaw)

	lines := []string{
		"2024-01-01T00:00:00Z INFO starting",
		"2024-01-01T00:00:01Z DEBUG ready",
	}
	for _, l := range lines {
		if err := w.WriteLine(l); err != nil {
			t.Fatalf("WriteLine: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	got := buf.String()
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("expected output to contain %q", l)
		}
	}
	if w.Count() != 2 {
		t.Errorf("Count() = %d, want 2", w.Count())
	}
}

func TestWriter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatJSON)

	if err := w.WriteLine("hello world"); err != nil {
		t.Fatalf("WriteLine: %v", err)
	}
	_ = w.Flush()

	got := buf.String()
	if !strings.Contains(got, `"line"`) {
		t.Errorf("JSON output missing \"line\" key, got: %s", got)
	}
	if !strings.Contains(got, "hello world") {
		t.Errorf("JSON output missing original content, got: %s", got)
	}
}

func TestWriter_Count(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatRaw)

	for i := 0; i < 5; i++ {
		_ = w.WriteLine("line")
	}
	if w.Count() != 5 {
		t.Errorf("Count() = %d, want 5", w.Count())
	}
}
