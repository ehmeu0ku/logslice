// Package tail provides functionality for following log files in real-time,
// similar to `tail -f`, emitting new lines as they are appended.
package tail

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourusername/logslice/internal/parser"
)

// FollowOptions configures the behaviour of a follow session.
type FollowOptions struct {
	// PollInterval controls how frequently the file is polled for new content
	// when no inotify/kqueue support is available. Defaults to 250ms.
	PollInterval time.Duration

	// StartAtEnd, when true, skips existing content and only emits lines
	// appended after the follow session begins.
	StartAtEnd bool
}

// DefaultFollowOptions returns a FollowOptions with sensible defaults.
func DefaultFollowOptions() FollowOptions {
	return FollowOptions{
		PollInterval: 250 * time.Millisecond,
		StartAtEnd:   false,
	}
}

// Follow opens the file at path and streams newly appended lines to the
// returned channel until ctx is cancelled or an unrecoverable error occurs.
// Each emitted *parser.LogLine has its timestamp parsed where possible.
// The returned error channel receives at most one error before being closed.
func Follow(ctx context.Context, path string, opts FollowOptions) (<-chan *parser.LogLine, <-chan error) {
	lines := make(chan *parser.LogLine, 64)
	errs := make(chan error, 1)

	go func() {
		defer close(lines)
		defer close(errs)

		f, err := os.Open(path)
		if err != nil {
			errs <- fmt.Errorf("follow: open %q: %w", path, err)
			return
		}
		defer f.Close()

		// Optionally seek to the end so we only see new content.
		if opts.StartAtEnd {
			if _, err := f.Seek(0, io.SeekEnd); err != nil {
				errs <- fmt.Errorf("follow: seek %q: %w", path, err)
				return
			}
		}

		reader := bufio.NewReader(f)
		pollInterval := opts.PollInterval
		if pollInterval <= 0 {
			pollInterval = DefaultFollowOptions().PollInterval
		}

		for {
			// Drain all currently available lines before sleeping.
			for {
				raw, readErr := reader.ReadString('\n')
				if len(raw) > 0 {
					// Strip trailing newline for consistent handling.
					if len(raw) > 0 && raw[len(raw)-1] == '\n' {
						raw = raw[:len(raw)-1]
					}
					line := parser.ParseLine(raw)
					select {
					case lines <- line:
					case <-ctx.Done():
						return
					}
				}
				if readErr != nil {
					// io.EOF is expected when we've caught up; any other error is fatal.
					if readErr != io.EOF {
						errs <- fmt.Errorf("follow: read %q: %w", path, readErr)
						return
					}
					break
				}
			}

			// Wait before polling again, respecting context cancellation.
			select {
			case <-ctx.Done():
				return
			case <-time.After(pollInterval):
			}
		}
	}()

	return lines, errs
}
