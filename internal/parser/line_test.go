package parser

import (
	"testing"
	"time"
)

func TestParseLine_WithTimestamp(t *testing.T) {
	raw := "2024-03-15T10:22:33Z INFO server started"
	l := ParseLine(raw)

	if !l.HasTime {
		t.Fatal("expected HasTime to be true")
	}
	if l.Timestamp.Year() != 2024 {
		t.Errorf("year: got %d, want 2024", l.Timestamp.Year())
	}
	if l.Raw != raw {
		t.Errorf("Raw mismatch: got %q", l.Raw)
	}
}

func TestParseLine_WithoutTimestamp(t *testing.T) {
	raw := "INFO no timestamp here"
	l := ParseLine(raw)

	if l.HasTime {
		t.Error("expected HasTime to be false")
	}
	if l.Raw != raw {
		t.Errorf("Raw mismatch: got %q", l.Raw)
	}
}

func TestLogLine_InRange(t *testing.T) {
	base := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	start := base
	end := base.Add(1 * time.Hour)

	tests := []struct {
		name      string
		lineTime  time.Time
		hasTime   bool
		wantRange bool
	}{
		{"exactly at start", start, true, true},
		{"exactly at end", end, true, true},
		{"in middle", base.Add(30 * time.Minute), true, true},
		{"before start", base.Add(-1 * time.Second), true, false},
		{"after end", end.Add(1 * time.Second), true, false},
		{"no timestamp", time.Time{}, false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := LogLine{Raw: "test", Timestamp: tc.lineTime, HasTime: tc.hasTime}
			if got := l.InRange(start, end); got != tc.wantRange {
				t.Errorf("InRange = %v, want %v", got, tc.wantRange)
			}
		})
	}
}
