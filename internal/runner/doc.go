// Package runner wires together the filter, parser, and output components
// into a single execution pipeline for the logslice tool.
//
// The Runner accepts a Config, opens the specified input file, applies
// timestamp range and optional pattern filters, and writes matching log
// lines to the configured output in the requested format.
//
// Typical usage:
//
//	cfg, err := config.ParseFlags()
//	if err != nil {
//		log.Fatal(err)
//	}
//	r := runner.New(cfg)
//	count, err := r.Run()
package runner
