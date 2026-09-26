package report

import (
	"github.com/c4rb0nx1/tuprwre/internal/event"
	"github.com/c4rb0nx1/tuprwre/internal/gateway"
)

// Redact applies the default gateway redactor to each event's argv,
// arguments and result, returning the number of replacements made. It is the
// read-side counterpart of the write-time redaction: a log written with
// redaction disabled, or produced by a third-party sensor, may still hold
// secrets, and the summary the report prints (and dumps as --json) is derived
// from exactly those fields. Redacting before Build keeps a secret out of the
// report regardless of how the log was produced.
//
// The redactor rewrites only secret-looking substrings and the values of
// credential-named fields; the rule inputs the report works from (commands,
// branch and context names, file paths) are unaffected.
func Redact(events []event.Event) int {
	count := 0
	for i := range events {
		before := events[i].RedactionCount
		events[i] = gateway.DefaultRedactor(events[i])
		count += int(events[i].RedactionCount - before)
	}
	return count
}
