package filter

import (
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// RangeFilter holds the start and end time boundaries for log filtering.
type RangeFilter struct {
	Start *time.Time
	End   *time.Time
}

// NewRangeFilter creates a RangeFilter from optional start and end timestamp strings.
// Empty strings are treated as unbounded.
func NewRangeFilter(start, end string) (*RangeFilter, error) {
	rf := &RangeFilter{}

	if start != "" {
		t, err := parser.ParseTimestamp(start)
		if err != nil {
			return nil, fmt.Errorf("invalid start timestamp %q: %w", start, err)
		}
		rf.Start = &t
	}

	if end != "" {
		t, err := parser.ParseTimestamp(end)
		if err != nil {
			return nil, fmt.Errorf("invalid end timestamp %q: %w", end, err)
		}
		rf.End = &t
	}

	return rf, nil
}

// Match returns true if the given log line falls within the filter's time range.
// Lines without a parsed timestamp are excluded when any boundary is set.
func (rf *RangeFilter) Match(line parser.LogLine) bool {
	if rf.Start == nil && rf.End == nil {
		return true
	}
	if line.Timestamp == nil {
		return false
	}
	return line.InRange(rf.Start, rf.End)
}
