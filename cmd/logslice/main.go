package main

import (
	"fmt"
	"os"

	"github.com/logslice/logslice/internal/config"
	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/runner"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: invalid config: %v\n", err)
		os.Exit(1)
	}

	r := runner.New(cfg)
	count, err := r.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}

	summary := &output.Summary{
		Matched: count,
		Input:   cfg.InputFile,
	}
	summary.Print(os.Stderr)
}
