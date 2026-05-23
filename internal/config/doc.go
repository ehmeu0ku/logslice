// Package config provides runtime configuration types and flag parsing
// for the logslice tool.
//
// Usage:
//
//	cfg, err := config.ParseFlags(os.Args[1:])
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := cfg.Validate(); err != nil {
//		log.Fatal(err)
//	}
//
// The Config struct holds all settings needed to drive a logslice run,
// including input/output paths, time range bounds, regex pattern,
// output format, and summary toggle.
package config
