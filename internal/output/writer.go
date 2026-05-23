package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Format defines the output format for log lines.
type Format string

const (
	FormatRaw  Format = "raw"
	FormatJSON Format = "json"
)

// Writer writes filtered log lines to a destination.
type Writer struct {
	w      *bufio.Writer
	format Format
	count  int
}

// NewWriter creates a new Writer targeting the given io.Writer.
func NewWriter(w io.Writer, format Format) *Writer {
	return &Writer{
		w:      bufio.NewWriter(w),
		format: format,
	}
}

// NewFileWriter opens a file for writing and returns a Writer.
func NewFileWriter(path string, format Format) (*Writer, *os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("output: create file %q: %w", path, err)
	}
	return NewWriter(f, format), f, nil
}

// WriteLine writes a single log line according to the configured format.
func (w *Writer) WriteLine(line string) error {
	var out string
	switch w.format {
	case FormatJSON:
		out = fmt.Sprintf(`{"line":%q}\n`, line)
	default:
		out = line + "\n"
	}
	if _, err := w.w.WriteString(out); err != nil {
		return fmt.Errorf("output: write line: %w", err)
	}
	w.count++
	return nil
}

// Flush flushes any buffered data to the underlying writer.
func (w *Writer) Flush() error {
	return w.w.Flush()
}

// Count returns the number of lines written.
func (w *Writer) Count() int {
	return w.count
}
