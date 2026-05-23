package filter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

const sampleLog = `2024-01-15T10:00:00Z INFO  starting server
2024-01-15T10:05:00Z DEBUG received request
2024-01-15T10:10:00Z ERROR connection refused
2024-01-15T10:15:00Z INFO  shutting down
no-timestamp line should be excluded when filter is set
`

func TestScanner_FullRange(t *testing.T) {
	rf := &filter.RangeFilter{}
	s := filter.NewScanner(strings.NewReader(sampleLog), rf)
	lines, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}
}

func TestScanner_StartBound(t *testing.T) {
	start := mustTime("2024-01-15T10:08:00Z")
	rf := &filter.RangeFilter{Start: &start}
	s := filter.NewScanner(strings.NewReader(sampleLog), rf)
	lines, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestScanner_StartAndEnd(t *testing.T) {
	start := mustTime("2024-01-15T10:04:00Z")
	end := mustTime("2024-01-15T10:11:00Z")
	rf := &filter.RangeFilter{Start: &start, End: &end}
	s := filter.NewScanner(strings.NewReader(sampleLog), rf)
	lines, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Errorf("expected 2 lines (10:05 and 10:10), got %d", len(lines))
	}
}

func TestScanner_NoMatch(t *testing.T) {
	start := mustTime("2024-01-15T12:00:00Z")
	rf := &filter.RangeFilter{Start: &start}
	s := filter.NewScanner(strings.NewReader(sampleLog), rf)
	lines, err := s.Scan()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(lines))
	}
}
