package glog

import (
	"testing"

	"go.opencensus.io/trace"
)

func Test_fromSpanContext(t *testing.T) {
	ProjectID = "my-projectid"

	tests := []struct {
		name   string
		span   trace.SpanContext
		trace  string
		spanID string
	}{
		{
			"span",
			trace.SpanContext{
				TraceID: [16]byte{0x01},
				SpanID:  [8]byte{0x02},
			},
			"projects/my-projectid/traces/01000000000000000000000000000000",
			"0200000000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trace, spanID := fromSpanContext(tt.span)
			if trace != tt.trace {
				t.Errorf("fromSpanContext() trace = %q, want %q", trace, tt.trace)
			}
			if spanID != tt.spanID {
				t.Errorf("fromSpanContext() spanID = %q, want %q", spanID, tt.spanID)
			}
		})
	}
}

func Test_parseTraceContext(t *testing.T) {
	ProjectID = "my-projectid"

	tests := []struct {
		name   string
		header string
		trace  string
		spanID string
	}{
		{"no header", "", "", ""},
		{"no span", "06796866738c859f2f19b7cfb3214824/0;o=1", "projects/my-projectid/traces/06796866738c859f2f19b7cfb3214824", ""},
		{"hex span", "06796866738c859f2f19b7cfb3214824/74;o=1", "projects/my-projectid/traces/06796866738c859f2f19b7cfb3214824", "000000000000004a"},
		{"with span", "06796866738c859f2f19b7cfb3214824/1;o=1", "projects/my-projectid/traces/06796866738c859f2f19b7cfb3214824", "0000000000000001"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trace, spanID := parseTraceContext(tt.header)
			if trace != tt.trace {
				t.Errorf("parseTraceContext() trace = %q, want %q", trace, tt.trace)
			}
			if spanID != tt.spanID {
				t.Errorf("parseTraceContext() spanID = %q, want %q", spanID, tt.spanID)
			}
		})
	}
}
