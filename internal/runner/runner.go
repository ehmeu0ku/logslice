package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/logslice/logslice/internal/config"
	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/output"
)

// Runner orchestrates the full log slicing pipeline.
type Runner struct {
	cfg    *config.Config
	stdout io.Writer
}

// New creates a new Runner with the given config.
func New(cfg *config.Config) *Runner {
	return &Runner{cfg: cfg, stdout: os.Stdout}
}

// NewWithWriter creates a Runner that writes output to the provided writer.
func NewWithWriter(cfg *config.Config, w io.Writer) *Runner {
	return &Runner{cfg: cfg, stdout: w}
}

// Run executes the log slicing pipeline and returns the number of matched lines.
func (r *Runner) Run() (int, error) {
	input, err := os.Open(r.cfg.InputFile)
	if err != nil {
		return 0, fmt.Errorf("opening input file: %w", err)
	}
	defer input.Close()

	rf := filter.NewRangeFilter(r.cfg.Start, r.cfg.End)

	var pf *filter.PatternFilter
	if r.cfg.Pattern != "" {
		pf, err = filter.NewPatternFilter(r.cfg.Pattern)
		if err != nil {
			return 0, fmt.Errorf("compiling pattern: %w", err)
		}
	}

	writer := output.NewWriter(r.stdout, r.cfg.Format)
	scanner := filter.NewScanner(input, rf, pf)

	count := 0
	for scanner.Scan() {
		line := scanner.Line()
		if err := writer.Write(line); err != nil {
			return count, fmt.Errorf("writing line: %w", err)
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		return count, fmt.Errorf("scanning input: %w", err)
	}

	return count, nil
}
