package output

import (
	"fmt"
	"io"
	"time"
)

// Summary holds statistics about a completed filter run.
type Summary struct {
	LinesScanned int
	LinesMatched int
	Duration     time.Duration
	OutputFormat Format
}

// Print writes a human-readable summary to w.
func (s Summary) Print(w io.Writer) {
	fmt.Fprintf(w, "--- logslice summary ---\n")
	fmt.Fprintf(w, "  scanned : %d lines\n", s.LinesScanned)
	fmt.Fprintf(w, "  matched : %d lines\n", s.LinesMatched)
	fmt.Fprintf(w, "  format  : %s\n", s.OutputFormat)
	fmt.Fprintf(w, "  elapsed : %s\n", s.Duration.Round(time.Millisecond))
}

// MatchRate returns the percentage of lines that matched, 0 if none scanned.
func (s Summary) MatchRate() float64 {
	if s.LinesScanned == 0 {
		return 0
	}
	return float64(s.LinesMatched) / float64(s.LinesScanned) * 100
}
