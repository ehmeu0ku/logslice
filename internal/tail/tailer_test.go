package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/tail"
)

func writeTempLog(t *testing.T, lines ...string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f
}

func TestTailer_EmitsNewLines(t *testing.T) {
	f := writeTempLog(t) // empty to start

	tlr := tail.NewTailer(f.Name(), 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lines, errs := tlr.Tail(ctx)

	// Write new content after tailing has started.
	go func() {
		time.Sleep(50 * time.Millisecond)
		f.WriteString("2024-01-01T00:00:01Z INFO hello\n")
		f.WriteString("2024-01-01T00:00:02Z INFO world\n")
	}()

	var got []string
	for line := range lines {
		got = append(got, line)
		if len(got) == 2 {
			cancel()
		}
	}

	select {
	case err := <-errs:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	default:
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestTailer_MissingFile(t *testing.T) {
	tlr := tail.NewTailer("/nonexistent/path/file.log", 10*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, errs := tlr.Tail(ctx)

	select {
	case err := <-errs:
		if err == nil {
			t.Fatal("expected an error for missing file, got nil")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for error")
	}
}

func TestTailer_ContextCancel(t *testing.T) {
	f := writeTempLog(t)
	tlr := tail.NewTailer(f.Name(), 10*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	lines, _ := tlr.Tail(ctx)
	cancel()

	// Channel should close promptly after cancel.
	select {
	case _, ok := <-lines:
		_ = ok // drained or closed — both acceptable
	case <-time.After(500 * time.Millisecond):
		t.Fatal("channel did not close after context cancel")
	}
}
