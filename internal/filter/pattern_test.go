package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func TestNewPatternFilter_ValidPattern(t *testing.T) {
	f, err := NewPatternFilter(`ERROR`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNewPatternFilter_InvalidPattern(t *testing.T) {
	_, err := NewPatternFilter(`[invalid`)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestPatternFilter_Match(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		raw     string
		want    bool
	}{
		{"exact word match", `ERROR`, "2024-01-01 ERROR something failed", true},
		{"no match", `ERROR`, "2024-01-01 INFO all good", false},
		{"regex group match", `(WARN|ERROR)`, "2024-01-01 WARN disk low", true},
		{"case sensitive no match", `error`, "2024-01-01 ERROR boom", false},
		{"partial match", `fail`, "2024-01-01 ERROR connection failure", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewPatternFilter(tt.pattern)
			if err != nil {
				t.Fatalf("NewPatternFilter(%q): %v", tt.pattern, err)
			}
			line := parser.LogLine{Raw: tt.raw}
			got := f.Match(line)
			if got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}

func TestPatternFilter_MatchString(t *testing.T) {
	f, _ := NewPatternFilter(`\d{3}`)
	if !f.MatchString("error code 404 returned") {
		t.Error("expected match for string containing digits")
	}
	if f.MatchString("no digits here") {
		t.Error("expected no match for string without digits")
	}
}

func TestPatternFilter_Pattern(t *testing.T) {
	pattern := `^INFO.*timeout`
	f, _ := NewPatternFilter(pattern)
	if f.Pattern() != pattern {
		t.Errorf("Pattern() = %q, want %q", f.Pattern(), pattern)
	}
}
