// Package tail implements live log file tailing for logslice.
//
// It provides a Tailer type that watches a file for new content and
// streams newly appended lines over a channel. This is useful for
// monitoring log files in real time while still applying logslice's
// timestamp and pattern filters.
//
// Usage:
//
//	tlr := tail.NewTailer("/var/log/app.log", 100*time.Millisecond)
//	lines, errs := tlr.Tail(ctx)
//	for line := range lines {
//		fmt.Println(line)
//	}
//	if err := <-errs; err != nil {
//		log.Fatal(err)
//	}
package tail
