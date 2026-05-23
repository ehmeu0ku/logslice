package filter

import (
	"regexp"

	"github.com/user/logslice/internal/parser"
)

// PatternFilter filters log lines by a regular expression pattern.
type PatternFilter struct {
	re *regexp.Regexp
}

// NewPatternFilter compiles the given pattern and returns a PatternFilter.
// Returns an error if the pattern is not a valid regular expression.
func NewPatternFilter(pattern string) (*PatternFilter, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &PatternFilter{re: re}, nil
}

// Match returns true if the log line's raw message matches the compiled pattern.
func (f *PatternFilter) Match(line parser.LogLine) bool {
	return f.re.MatchString(line.Raw)
}

// MatchString returns true if the given raw string matches the compiled pattern.
func (f *PatternFilter) MatchString(s string) bool {
	return f.re.MatchString(s)
}

// Pattern returns the original pattern string used to build the filter.
func (f *PatternFilter) Pattern() string {
	return f.re.String()
}
