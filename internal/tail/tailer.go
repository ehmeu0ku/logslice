// Package tail provides live log tailing functionality for logslice.
package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Tailer watches a file and emits new lines as they are written.
type Tailer struct {
	path     string
	pollRate time.Duration
}

// NewTailer creates a Tailer for the given file path.
// pollRate controls how often the file is checked for new content.
func NewTailer(path string, pollRate time.Duration) *Tailer {
	return &Tailer{
		path:     path,
		pollRate: pollRate,
	}
}

// Tail reads new lines from the file and sends them to the returned channel.
// It seeks to the end of the file before starting. The channel is closed when
// ctx is cancelled or a read error occurs.
func (t *Tailer) Tail(ctx context.Context) (<-chan string, <-chan error) {
	lines := make(chan string, 64)
	errs := make(chan error, 1)

	go func() {
		defer close(lines)
		defer close(errs)

		f, err := os.Open(t.path)
		if err != nil {
			errs <- err
			return
		}
		defer f.Close()

		// Seek to end so we only emit new lines.
		if _, err := f.Seek(0, io.SeekEnd); err != nil {
			errs <- err
			return
		}

		scanner := bufio.NewScanner(f)
		ticker := time.NewTicker(t.pollRate)
		defer ticker.Stop()

		for {
			for scanner.Scan() {
				select {
				case lines <- scanner.Text():
				case <-ctx.Done():
					return
				}
			}
			if err := scanner.Err(); err != nil {
				errs <- err
				return
			}
			select {
				case <-ticker.C:
				case <-ctx.Done():
					return
			}
		}
	}()

	return lines, errs
}
