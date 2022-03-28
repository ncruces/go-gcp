// Package glog implements leveled, structured logging for Google App Engine,
// Kubernetes Engine, Cloud Run, and Cloud Functions.
//
// Structured logs are written to stdout or stderr as JSON objects
// serialized on a single line.
// The Logging agent then sends the structured logs to Cloud Logging
// as the jsonPayload of the LogEntry structure.
package glog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/functions/metadata"
	"go.opencensus.io/trace"
)

var std Logger = Logger{callers: 1}

// ProjectID should be set to the Google Cloud project ID.
var ProjectID string = os.Getenv("GOOGLE_CLOUD_PROJECT")

// LogSourceLocation should be set to false to avoid associating
// source code location information with the entry.
var LogSourceLocation bool = true

// Print logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...any) {
	std.Print(v...)
}

// Println logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...any) {
	std.Println(v...)
}

// Printf logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...any) {
	std.Printf(format, v...)
}

// Printj logs an entry with no assigned severity level.
// Arguments populate jsonPayload in the log entry.
func Printj(msg string, v any) {
	std.Printj(msg, v)
}

// Printw logs an entry with no assigned severity level.
// Arguments populate jsonPayload in the log entry.
func Printw(msg string, kvs ...any) {
	std.Printw(msg, kvs...)
}

// Debug logs debug or trace information.
// Arguments are handled in the manner of fmt.Print.
func Debug(v ...any) {
	std.Debug(v...)
}

// Debugln logs debug or trace information.
// Arguments are handled in the manner of fmt.Println.
func Debugln(v ...any) {
	std.Debugln(v...)
}

// Debugf logs debug or trace information.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...any) {
	std.Debugf(format, v...)
}

// Debugj logs debug or trace information.
// Arguments populate jsonPayload in the log entry.
func Debugj(msg string, v any) {
	std.Debugj(msg, v)
}

// Debugw logs debug or trace information.
// Arguments populate jsonPayload in the log entry.
func Debugw(msg string, kvs ...any) {
	std.Debugw(msg, kvs...)
}

// Info logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...any) {
	std.Info(v...)
}

// Infoln logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Println.
func Infoln(v ...any) {
	std.Infoln(v...)
}

// Infof logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...any) {
	std.Infof(format, v...)
}

// Infoj logs routine information, such as ongoing status or performance.
// Arguments populate jsonPayload in the log entry.
func Infoj(msg string, v any) {
	std.Infoj(msg, v)
}

// Infow logs routine information, such as ongoing status or performance.
// Arguments populate jsonPayload in the log entry.
func Infow(msg string, kvs ...any) {
	std.Infow(msg, kvs...)
}

// Notice logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Print.
func Notice(v ...any) {
	std.Notice(v...)
}

// Noticeln logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Println.
func Noticeln(v ...any) {
	std.Noticeln(v...)
}

// Noticef logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Printf.
func Noticef(format string, v ...any) {
	std.Noticef(format, v...)
}

// Noticej logs normal but significant events, such as start up, shut down, or configuration.
// Arguments populate jsonPayload in the log entry.
func Noticej(msg string, v any) {
	std.Noticej(msg, v)
}

// Noticew logs normal but significant events, such as start up, shut down, or configuration.
// Arguments populate jsonPayload in the log entry.
func Noticew(msg string, kvs ...any) {
	std.Noticew(msg, kvs...)
}

// Warning logs events that might cause problems.
// Arguments are handled in the manner of fmt.Print.
func Warning(v ...any) {
	std.Warning(v...)
}

// Warningln logs events that might cause problems.
// Arguments are handled in the manner of fmt.Println.
func Warningln(v ...any) {
	std.Warningln(v...)
}

// Warningf logs events that might cause problems.
// Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, v ...any) {
	std.Warningf(format, v...)
}

// Warningj logs events that might cause problems.
// Arguments populate jsonPayload in the log entry.
func Warningj(msg string, v any) {
	std.Warningj(msg, v)
}

// Warningw logs events that might cause problems.
// Arguments populate jsonPayload in the log entry.
func Warningw(msg string, kvs ...any) {
	std.Warningw(msg, kvs...)
}

// Error logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...any) {
	std.Error(v...)
}

// Errorln logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Println.
func Errorln(v ...any) {
	std.Errorln(v...)
}

// Errorf logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...any) {
	std.Errorf(format, v...)
}

// Errorj logs events likely to cause problems.
// Arguments populate jsonPayload in the log entry.
func Errorj(msg string, v any) {
	std.Errorj(msg, v)
}

// Errorw logs events likely to cause problems.
// Arguments populate jsonPayload in the log entry.
func Errorw(msg string, kvs ...any) {
	std.Errorw(msg, kvs...)
}

// Critical logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Print.
func Critical(v ...any) {
	std.Critical(v...)
}

// Criticalln logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Println.
func Criticalln(v ...any) {
	std.Criticalln(v...)
}

// Criticalf logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Printf.
func Criticalf(format string, v ...any) {
	std.Criticalf(format, v...)
}

// Criticalj logs events that cause more severe problems or outages.
// Arguments populate jsonPayload in the log entry.
func Criticalj(msg string, v any) {
	std.Criticalj(msg, v)
}

// Criticalw logs events that cause more severe problems or outages.
// Arguments populate jsonPayload in the log entry.
func Criticalw(msg string, kvs ...any) {
	std.Criticalw(msg, kvs...)
}

// Alert logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Print.
func Alert(v ...any) {
	std.Alert(v...)
}

// Alertln logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Println.
func Alertln(v ...any) {
	std.Alertln(v...)
}

// Alertf logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Printf.
func Alertf(format string, v ...any) {
	std.Alertf(format, v...)
}

// Alertj logs when a person must take an action immediately.
// Arguments populate jsonPayload in the log entry.
func Alertj(msg string, v any) {
	std.Alertj(msg, v)
}

// Alertw logs when a person must take an action immediately.
// Arguments populate jsonPayload in the log entry.
func Alertw(msg string, kvs ...any) {
	std.Alertw(msg, kvs...)
}

// Emergency logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Print.
func Emergency(v ...any) {
	std.Emergency(v...)
}

// Emergencyln logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Println.
func Emergencyln(v ...any) {
	std.Emergencyln(v...)
}

// Emergencyf logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Printf.
func Emergencyf(format string, v ...any) {
	std.Emergencyf(format, v...)
}

// Emergencyj logs when one or more systems are unusable.
// Arguments populate jsonPayload in the log entry.
func Emergencyj(msg string, v any) {
	std.Emergencyj(msg, v)
}

// Emergencyw logs when one or more systems are unusable.
// Arguments populate jsonPayload in the log entry.
func Emergencyw(msg string, kvs ...any) {
	std.Emergencyw(msg, kvs...)
}

// A Logger that logs entries with additional context.
type Logger struct {
	callers     int
	trace       string
	spanID      string
	executionID string
	request     *httpRequest
}

// ForRequest creates a Logger with metadata from an http.Request.
func ForRequest(r *http.Request) (l Logger) {
	l.trace, l.spanID = parseTraceContext(r.Header.Get("X-Cloud-Trace-Context"))
	l.executionID = r.Header.Get("Function-Execution-Id")
	l.request = &httpRequest{
		RequestMethod: r.Method,
		RequestUrl:    r.RequestURI,
		UserAgent:     r.UserAgent(),
		RemoteIp:      r.RemoteAddr,
		Referer:       r.Referer(),
		Protocol:      r.Proto,
	}
	return l
}

// ForContext creates a Logger with metadata from a context.Context.
func ForContext(ctx context.Context) (l Logger) {
	l.SetContext(ctx)
	return l
}

// SetContext updates a Logger with metadata from a context.Context.
func (l *Logger) SetContext(ctx context.Context) {
	if span := trace.FromContext(ctx); span != nil {
		l.trace, l.spanID = fromSpanContext(span.SpanContext())
	}
	if meta, _ := metadata.FromContext(ctx); meta != nil {
		l.executionID = meta.EventID
	}
}

// Print logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Print(v ...any) {
	logm(defaultsv, l, v...)
}

// Println logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Println(v ...any) {
	logn(defaultsv, l, v...)
}

// Printf logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Printf(format string, v ...any) {
	logf(defaultsv, l, format, v...)
}

// Printj logs an entry with no assigned severity level.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Printj(msg string, v any) {
	logj(defaultsv, l, msg, v)
}

// Printw logs an entry with no assigned severity level.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Printw(msg string, kvs ...any) {
	logw(defaultsv, l, msg, kvs)
}

// Debug logs debug or trace information.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Debug(v ...any) {
	logm(debugsv, l, v...)
}

// Debugln logs debug or trace information.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Debugln(v ...any) {
	logn(debugsv, l, v...)
}

// Debugf logs debug or trace information.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Debugf(format string, v ...any) {
	logf(debugsv, l, format, v...)
}

// Debugj logs debug or trace information.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Debugj(msg string, v any) {
	logj(debugsv, l, msg, v)
}

// Debugw logs debug or trace information.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Debugw(msg string, kvs ...any) {
	logw(debugsv, l, msg, kvs)
}

// Info logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Info(v ...any) {
	logm(infosv, l, v...)
}

// Infoln logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Infoln(v ...any) {
	logn(infosv, l, v...)
}

// Infof logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Infof(format string, v ...any) {
	logf(infosv, l, format, v...)
}

// Infoj logs routine information, such as ongoing status or performance.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Infoj(msg string, v any) {
	logj(infosv, l, msg, v)
}

// Infow logs routine information, such as ongoing status or performance.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Infow(msg string, kvs ...any) {
	logw(infosv, l, msg, kvs)
}

// Notice logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Notice(v ...any) {
	logm(noticesv, l, v...)
}

// Noticeln logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Noticeln(v ...any) {
	logn(noticesv, l, v...)
}

// Noticef logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Noticef(format string, v ...any) {
	logf(noticesv, l, format, v...)
}

// Noticej logs normal but significant events, such as start up, shut down, or configuration.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Noticej(msg string, v any) {
	logj(noticesv, l, msg, v)
}

// Noticew logs normal but significant events, such as start up, shut down, or configuration.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Noticew(msg string, kvs ...any) {
	logw(noticesv, l, msg, kvs)
}

// Warning logs events that might cause problems.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Warning(v ...any) {
	logm(warningsv, l, v...)
}

// Warningln logs events that might cause problems.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Warningln(v ...any) {
	logn(warningsv, l, v...)
}

// Warningf logs events that might cause problems.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Warningf(format string, v ...any) {
	logf(warningsv, l, format, v...)
}

// Warningj logs events that might cause problems.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Warningj(msg string, v any) {
	logj(warningsv, l, msg, v)
}

// Warningw logs events that might cause problems.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Warningw(msg string, kvs ...any) {
	logw(warningsv, l, msg, kvs)
}

// Error logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Error(v ...any) {
	logm(errorsv, l, v...)
}

// Errorln logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Errorln(v ...any) {
	logn(errorsv, l, v...)
}

// Errorf logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Errorf(format string, v ...any) {
	logf(errorsv, l, format, v...)
}

// Errorj logs events likely to cause problems.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Errorj(msg string, v any) {
	logj(errorsv, l, msg, v)
}

// Errorw logs events likely to cause problems.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Errorw(msg string, kvs ...any) {
	logw(errorsv, l, msg, kvs)
}

// Critical logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Critical(v ...any) {
	logm(criticalsv, l, v...)
}

// Criticalln logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Criticalln(v ...any) {
	logn(criticalsv, l, v...)
}

// Criticalf logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Criticalf(format string, v ...any) {
	logf(criticalsv, l, format, v...)
}

// Criticalj logs events that cause more severe problems or outages.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Criticalj(msg string, v any) {
	logj(criticalsv, l, msg, v)
}

// Criticalw logs events that cause more severe problems or outages.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Criticalw(msg string, kvs ...any) {
	logw(criticalsv, l, msg, kvs)
}

// Alert logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Alert(v ...any) {
	logm(alertsv, l, v...)
}

// Alertln logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Alertln(v ...any) {
	logn(alertsv, l, v...)
}

// Alertf logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Alertf(format string, v ...any) {
	logf(alertsv, l, format, v...)
}

// Alertj logs when a person must take an action immediately.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Alertj(msg string, v any) {
	logj(alertsv, l, msg, v)
}

// Alertw logs when a person must take an action immediately.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Alertw(msg string, kvs ...any) {
	logw(alertsv, l, msg, kvs)
}

// Emergency logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Emergency(v ...any) {
	logm(emergencysv, l, v...)
}

// Emergencyln logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Emergencyln(v ...any) {
	logn(emergencysv, l, v...)
}

// Emergencyf logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Emergencyf(format string, v ...any) {
	logf(emergencysv, l, format, v...)
}

// Emergencyj logs when one or more systems are unusable.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Emergencyj(msg string, v any) {
	logj(emergencysv, l, msg, v)
}

// Emergencyw logs when one or more systems are unusable.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Emergencyw(msg string, kvs ...any) {
	logw(emergencysv, l, msg, kvs)
}

type severity int32

const (
	defaultsv severity = iota * 100
	debugsv
	infosv
	noticesv
	warningsv
	errorsv
	criticalsv
	alertsv
	emergencysv
)

func (s severity) String() string {
	switch s {
	default:
		return ""
	case debugsv:
		return "DEBUG"
	case infosv:
		return "INFO"
	case noticesv:
		return "NOTICE"
	case warningsv:
		return "WARNING"
	case errorsv:
		return "ERROR"
	case criticalsv:
		return "CRITICAL"
	case alertsv:
		return "ALERT"
	case emergencysv:
		return "EMERGENCY"
	}
}

func (s severity) File() *os.File {
	if s >= errorsv {
		return os.Stderr
	} else {
		return os.Stdout
	}
}

func logm(s severity, l Logger, v ...any) {
	logs(s, l, fmt.Sprint(v...))
}

func logn(s severity, l Logger, v ...any) {
	logs(s, l, fmt.Sprintln(v...))
}

func logf(s severity, l Logger, format string, v ...any) {
	logs(s, l, fmt.Sprintf(format, v...))
}

func logs(s severity, l Logger, msg string) {
	entry := entry{
		Message:        strings.TrimSuffix(msg, "\n"),
		Severity:       s.String(),
		Trace:          l.trace,
		SpanID:         l.spanID,
		HttpRequest:    l.request,
		SourceLocation: location(4 + l.callers),
		Labels:         executionLabels(l.executionID),
	}
	json.NewEncoder(s.File()).Encode(entry)
}

func logj(s severity, l Logger, msg string, j any) {
	entry := make(map[string]json.RawMessage)
	if buf, err := json.Marshal(j); err != nil {
		panic(err)
	} else if err := json.Unmarshal(buf, &entry); err != nil {
		panic(err)
	}

	loge(s, l, msg, entry)
}

func logw(s severity, l Logger, msg string, kvs []any) {
	entry := make(map[string]json.RawMessage, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		var err error
		k, v := kvs[i].(string), kvs[i+1]
		entry[k], err = json.Marshal(v)
		if err != nil {
			panic(err)
		}
	}

	loge(s, l, msg, entry)
}

func loge(s severity, l Logger, msg string, entry map[string]json.RawMessage) {
	if v := msg; v != "" {
		entry["message"], _ = json.Marshal(v)
	}
	if v := s; v != 0 {
		entry["severity"], _ = json.Marshal(v.String())
	}
	if v := l.trace; v != "" {
		entry["logging.googleapis.com/trace"], _ = json.Marshal(v)
	}
	if v := l.spanID; v != "" {
		entry["logging.googleapis.com/spanId"], _ = json.Marshal(v)
	}
	if v := l.request; v != nil {
		entry["httpRequest"], _ = json.Marshal(v)
	}
	if v := l.executionID; v != "" {
		entry["labels"], _ = json.Marshal(executionLabels(l.executionID))
	}
	if v := location(4 + l.callers); v != nil {
		entry["logging.googleapis.com/sourceLocation"], _ = json.Marshal(v)
	}

	json.NewEncoder(s.File()).Encode(entry)
}

type entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`
	SpanID   string `json:"logging.googleapis.com/spanId,omitempty"`

	HttpRequest    *httpRequest    `json:"httpRequest,omitempty"`
	SourceLocation *sourceLocation `json:"logging.googleapis.com/sourceLocation,omitempty"`
	Labels         executionLabels `json:"logging.googleapis.com/labels,omitempty"`
}

type httpRequest struct {
	RequestMethod string `json:"requestMethod,omitempty"`
	RequestUrl    string `json:"requestUrl,omitempty"`
	UserAgent     string `json:"userAgent,omitempty"`
	RemoteIp      string `json:"remoteIp,omitempty"`
	Referer       string `json:"referer,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

type executionLabels string

func (e executionLabels) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ExecutionID string `json:"execution_id"`
	}{
		string(e),
	})
}

type sourceLocation struct {
	File     string `json:"file,omitempty"`
	Line     string `json:"line,omitempty"`
	Function string `json:"function,omitempty"`
}
