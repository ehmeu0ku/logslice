package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/output"
)

func TestSummary_Print(t *testing.T) {
	s := output.Summary{
		LinesScanned: 1000,
		LinesMatched: 42,
		Duration:     250 * time.Millisecond,
		OutputFormat: output.FormatRaw,
	}

	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	for _, want := range []string{"1000", "42", "raw", "250ms"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print() output missing %q\ngot: %s", want, out)
		}
	}
}

func TestSummary_MatchRate(t *testing.T) {
	tests := []struct {
		scanned, matched int
		want             float64
	}{
		{100, 50, 50.0},
		{200, 200, 100.0},
		{0, 0, 0.0},
		{1000, 1, 0.1},
	}
	for _, tc := range tests {
		s := output.Summary{LinesScanned: tc.scanned, LinesMatched: tc.matched}
		got := s.MatchRate()
		if got != tc.want {
			t.Errorf("MatchRate(%d/%d) = %.2f, want %.2f", tc.matched, tc.scanned, got, tc.want)
		}
	}
}
