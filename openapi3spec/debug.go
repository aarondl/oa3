package openapi3spec

import (
	"fmt"
	"os"
)

var (
	// DebugOutput controls whether or not debug output will be written
	DebugOutput = false
)

func debugln(intf ...interface{}) {
	if DebugOutput {
		_, _ = fmt.Fprintln(os.Stderr, intf...)
	}
}

func debugf(format string, intf ...interface{}) {
	if DebugOutput {
		_, _ = fmt.Fprintf(os.Stderr, format, intf...)
	}
}
