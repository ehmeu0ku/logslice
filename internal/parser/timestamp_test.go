package parser

import (
	"testing"
	"time"
)

func TestParseTimestamp_KnownFormats(t *testing.T) {
	cases := []struct {
		input    string
		wantYear int
		wantErr  bool
	}{
		{"2024-03-15T10:22:33Z", 2024, false},
		{"2024-03-15T10:22:33.123456789Z", 2024, false},
		{"2024-03-15 10:22:33", 2024, false},
		{"2024-03-15 10:22:33.000", 2024, false},
		{"15/Mar/2024:10:22:33 +0000", 2024, false},
		{"not-a-timestamp", 0, true},
		{"", 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, _, err := ParseTimestamp(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Year() != tc.wantYear {
				t.Errorf("year: got %d, want %d", got.Year(), tc.wantYear)
			}
		})
	}
}

func TestParseTimestampWithFormat(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		got, err := ParseTimestampWithFormat("2024-03-15 10:22:33", "2006-01-02 15:04:05")
		if err != nil {
			t.Fatal(err)
		}
		if got.Month() != time.March {
			t.Errorf("expected March, got %v", got.Month())
		}
	})

	t.Run("invalid", func(t *testing.T) {
		_, err := ParseTimestampWithFormat("bad", "2006-01-02")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestSupportedFormats(t *testing.T) {
	f := SupportedFormats()
	if len(f) == 0 {
		t.Error("expected at least one supported format")
	}
	// Mutation of returned slice must not affect internal list
	f[0] = "tampered"
	f2 := SupportedFormats()
	if f2[0] == "tampered" {
		t.Error("SupportedFormats should return a copy")
	}
}
