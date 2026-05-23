package filter

import (
	"bufio"
	"io"

	"github.com/yourorg/logslice/internal/parser"
)

// Scanner reads lines from a reader and emits those matching the RangeFilter.
type Scanner struct {
	reader  io.Reader
	filter  *RangeFilter
	results []parser.LogLine
}

// NewScanner creates a Scanner that applies the given RangeFilter.
func NewScanner(r io.Reader, rf *RangeFilter) *Scanner {
	return &Scanner{
		reader: r,
		filter: rf,
	}
}

// Scan reads all lines and collects those matching the filter.
// Returns the matched LogLine slice or an error.
func (s *Scanner) Scan() ([]parser.LogLine, error) {
	s.results = nil
	sc := bufio.NewScanner(s.reader)
	for sc.Scan() {
		raw := sc.Text()
		line := parser.ParseLine(raw)
		if s.filter.Match(line) {
			s.results = append(s.results, line)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return s.results, nil
}
