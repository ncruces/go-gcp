package glog

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"go.opencensus.io/trace"
)

func location(skip int) *sourceLocation {
	if !LogSourceLocation {
		return nil
	}
	if pc, file, line, ok := runtime.Caller(skip); ok {
		loc := &sourceLocation{
			File: file,
			Line: strconv.Itoa(line),
		}
		if f := runtime.FuncForPC(pc); f != nil {
			loc.Function = f.Name()
		}
		return loc
	}
	return nil
}

func fromSpanContext(spanContext trace.SpanContext) (trace, spanID string) {
	if ProjectID == "" {
		return
	}

	trace = fmt.Sprintf("projects/%s/traces/%s", ProjectID, spanContext.TraceID)
	spanID = spanContext.SpanID.String()
	return
}

func parseTraceContext(traceContext string) (trace, spanID string) {
	if traceContext == "" || ProjectID == "" {
		return
	}

	t, rest, ok := cut(traceContext, "/")
	if !ok {
		return
	}
	trace = fmt.Sprintf("projects/%s/traces/%s", ProjectID, t)

	s, _, ok := cut(rest, ";")
	if !ok {
		return
	}
	if s, _ := strconv.ParseUint(s, 10, 64); s > 0 {
		spanID = fmt.Sprintf("%016x", s)
	}

	return
}

// TODO: replace with strings.Cut.
func cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}
