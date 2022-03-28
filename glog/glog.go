// Package glog implements leveled, structured logging for Google App Engine,
// Cloud Run, and Cloud Functions.
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
func Print(v ...interface{}) {
	std.Print(v...)
}

// Println logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	std.Println(v...)
}

// Printf logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	std.Printf(format, v...)
}

// Printj logs an entry with no assigned severity level.
// Arguments populate jsonPayload in the log entry.
func Printj(msg string, v interface{}) {
	std.Printj(msg, v)
}

// Printw logs an entry with no assigned severity level.
// Arguments populate jsonPayload in the log entry.
func Printw(msg string, kvs ...interface{}) {
	std.Printw(msg, kvs...)
}

// Debug logs debug or trace information.
// Arguments are handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	std.Debug(v...)
}

// Debugln logs debug or trace information.
// Arguments are handled in the manner of fmt.Println.
func Debugln(v ...interface{}) {
	std.Debugln(v...)
}

// Debugf logs debug or trace information.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

// Debugj logs debug or trace information.
// Arguments populate jsonPayload in the log entry.
func Debugj(msg string, v interface{}) {
	std.Debugj(msg, v)
}

// Debugw logs debug or trace information.
// Arguments populate jsonPayload in the log entry.
func Debugw(msg string, kvs ...interface{}) {
	std.Debugw(msg, kvs...)
}

// Info logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	std.Info(v...)
}

// Infoln logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Println.
func Infoln(v ...interface{}) {
	std.Infoln(v...)
}

// Infof logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

// Infoj logs routine information, such as ongoing status or performance.
// Arguments populate jsonPayload in the log entry.
func Infoj(msg string, v interface{}) {
	std.Infoj(msg, v)
}

// Infow logs routine information, such as ongoing status or performance.
// Arguments populate jsonPayload in the log entry.
func Infow(msg string, kvs ...interface{}) {
	std.Infow(msg, kvs...)
}

// Notice logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Print.
func Notice(v ...interface{}) {
	std.Notice(v...)
}

// Noticeln logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Println.
func Noticeln(v ...interface{}) {
	std.Noticeln(v...)
}

// Noticef logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Printf.
func Noticef(format string, v ...interface{}) {
	std.Noticef(format, v...)
}

// Noticej logs normal but significant events, such as start up, shut down, or configuration.
// Arguments populate jsonPayload in the log entry.
func Noticej(msg string, v interface{}) {
	std.Noticej(msg, v)
}

// Noticew logs normal but significant events, such as start up, shut down, or configuration.
// Arguments populate jsonPayload in the log entry.
func Noticew(msg string, kvs ...interface{}) {
	std.Noticew(msg, kvs...)
}

// Warning logs events that might cause problems.
// Arguments are handled in the manner of fmt.Print.
func Warning(v ...interface{}) {
	std.Warning(v...)
}

// Warningln logs events that might cause problems.
// Arguments are handled in the manner of fmt.Println.
func Warningln(v ...interface{}) {
	std.Warningln(v...)
}

// Warningf logs events that might cause problems.
// Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, v ...interface{}) {
	std.Warningf(format, v...)
}

// Warningj logs events that might cause problems.
// Arguments populate jsonPayload in the log entry.
func Warningj(msg string, v interface{}) {
	std.Warningj(msg, v)
}

// Warningw logs events that might cause problems.
// Arguments populate jsonPayload in the log entry.
func Warningw(msg string, kvs ...interface{}) {
	std.Warningw(msg, kvs...)
}

// Error logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	std.Error(v...)
}

// Errorln logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	std.Errorln(v...)
}

// Errorf logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

// Errorj logs events likely to cause problems.
// Arguments populate jsonPayload in the log entry.
func Errorj(msg string, v interface{}) {
	std.Errorj(msg, v)
}

// Errorw logs events likely to cause problems.
// Arguments populate jsonPayload in the log entry.
func Errorw(msg string, kvs ...interface{}) {
	std.Errorw(msg, kvs...)
}

// Critical logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Print.
func Critical(v ...interface{}) {
	std.Critical(v...)
}

// Criticalln logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Println.
func Criticalln(v ...interface{}) {
	std.Criticalln(v...)
}

// Criticalf logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Printf.
func Criticalf(format string, v ...interface{}) {
	std.Criticalf(format, v...)
}

// Criticalj logs events that cause more severe problems or outages.
// Arguments populate jsonPayload in the log entry.
func Criticalj(msg string, v interface{}) {
	std.Criticalj(msg, v)
}

// Criticalw logs events that cause more severe problems or outages.
// Arguments populate jsonPayload in the log entry.
func Criticalw(msg string, kvs ...interface{}) {
	std.Criticalw(msg, kvs...)
}

// Alert logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Print.
func Alert(v ...interface{}) {
	std.Alert(v...)
}

// Alertln logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Println.
func Alertln(v ...interface{}) {
	std.Alertln(v...)
}

// Alertf logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Printf.
func Alertf(format string, v ...interface{}) {
	std.Alertf(format, v...)
}

// Alertj logs when a person must take an action immediately.
// Arguments populate jsonPayload in the log entry.
func Alertj(msg string, v interface{}) {
	std.Alertj(msg, v)
}

// Alertw logs when a person must take an action immediately.
// Arguments populate jsonPayload in the log entry.
func Alertw(msg string, kvs ...interface{}) {
	std.Alertw(msg, kvs...)
}

// Emergency logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Print.
func Emergency(v ...interface{}) {
	std.Emergency(v...)
}

// Emergencyln logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Println.
func Emergencyln(v ...interface{}) {
	std.Emergencyln(v...)
}

// Emergencyf logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Printf.
func Emergencyf(format string, v ...interface{}) {
	std.Emergencyf(format, v...)
}

// Emergencyj logs when one or more systems are unusable.
// Arguments populate jsonPayload in the log entry.
func Emergencyj(msg string, v interface{}) {
	std.Emergencyj(msg, v)
}

// Emergencyw logs when one or more systems are unusable.
// Arguments populate jsonPayload in the log entry.
func Emergencyw(msg string, kvs ...interface{}) {
	std.Emergencyw(msg, kvs...)
}

// A Logger that logs entries with additional context.
type Logger struct {
	trace       string
	spanID      string
	executionID string
	request     *httpRequest
	callers     int
}

// ForContext creates a Logger with metadata from a context.Context.
func ForContext(ctx context.Context) (l Logger) {
	if span := trace.FromContext(ctx); span != nil {
		l.trace, l.spanID = fromSpanContext(span.SpanContext())
	}
	if meta, _ := metadata.FromContext(ctx); meta != nil {
		l.executionID = meta.EventID
	}
	return l
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

// Print logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Print(v ...interface{}) {
	logm(defaultsv, l, v...)
}

// Println logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Println(v ...interface{}) {
	logn(defaultsv, l, v...)
}

// Printf logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Printf(format string, v ...interface{}) {
	logf(defaultsv, l, format, v...)
}

// Printj logs an entry with no assigned severity level.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Printj(msg string, v interface{}) {
	logj(defaultsv, l, msg, v)
}

// Printw logs an entry with no assigned severity level.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Printw(msg string, kvs ...interface{}) {
	logw(defaultsv, l, msg, kvs)
}

// Debug logs debug or trace information.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Debug(v ...interface{}) {
	logm(debugsv, l, v...)
}

// Debugln logs debug or trace information.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Debugln(v ...interface{}) {
	logn(debugsv, l, v...)
}

// Debugf logs debug or trace information.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Debugf(format string, v ...interface{}) {
	logf(debugsv, l, format, v...)
}

// Debugj logs debug or trace information.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Debugj(msg string, v interface{}) {
	logj(debugsv, l, msg, v)
}

// Debugw logs debug or trace information.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Debugw(msg string, kvs ...interface{}) {
	logw(debugsv, l, msg, kvs)
}

// Info logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Info(v ...interface{}) {
	logm(infosv, l, v...)
}

// Infoln logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Infoln(v ...interface{}) {
	logn(infosv, l, v...)
}

// Infof logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Infof(format string, v ...interface{}) {
	logf(infosv, l, format, v...)
}

// Infoj logs routine information, such as ongoing status or performance.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Infoj(msg string, v interface{}) {
	logj(infosv, l, msg, v)
}

// Infow logs routine information, such as ongoing status or performance.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Infow(msg string, kvs ...interface{}) {
	logw(infosv, l, msg, kvs)
}

// Notice logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Notice(v ...interface{}) {
	logm(noticesv, l, v...)
}

// Noticeln logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Noticeln(v ...interface{}) {
	logn(noticesv, l, v...)
}

// Noticef logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Noticef(format string, v ...interface{}) {
	logf(noticesv, l, format, v...)
}

// Noticej logs normal but significant events, such as start up, shut down, or configuration.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Noticej(msg string, v interface{}) {
	logj(noticesv, l, msg, v)
}

// Noticew logs normal but significant events, such as start up, shut down, or configuration.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Noticew(msg string, kvs ...interface{}) {
	logw(noticesv, l, msg, kvs)
}

// Warning logs events that might cause problems.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Warning(v ...interface{}) {
	logm(warningsv, l, v...)
}

// Warningln logs events that might cause problems.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Warningln(v ...interface{}) {
	logn(warningsv, l, v...)
}

// Warningf logs events that might cause problems.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Warningf(format string, v ...interface{}) {
	logf(warningsv, l, format, v...)
}

// Warningj logs events that might cause problems.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Warningj(msg string, v interface{}) {
	logj(warningsv, l, msg, v)
}

// Warningw logs events that might cause problems.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Warningw(msg string, kvs ...interface{}) {
	logw(warningsv, l, msg, kvs)
}

// Error logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Error(v ...interface{}) {
	logm(errorsv, l, v...)
}

// Errorln logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Errorln(v ...interface{}) {
	logn(errorsv, l, v...)
}

// Errorf logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Errorf(format string, v ...interface{}) {
	logf(errorsv, l, format, v...)
}

// Errorj logs events likely to cause problems.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Errorj(msg string, v interface{}) {
	logj(errorsv, l, msg, v)
}

// Errorw logs events likely to cause problems.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Errorw(msg string, kvs ...interface{}) {
	logw(errorsv, l, msg, kvs)
}

// Critical logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Critical(v ...interface{}) {
	logm(criticalsv, l, v...)
}

// Criticalln logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Criticalln(v ...interface{}) {
	logn(criticalsv, l, v...)
}

// Criticalf logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Criticalf(format string, v ...interface{}) {
	logf(criticalsv, l, format, v...)
}

// Criticalj logs events that cause more severe problems or outages.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Criticalj(msg string, v interface{}) {
	logj(criticalsv, l, msg, v)
}

// Criticalw logs events that cause more severe problems or outages.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Criticalw(msg string, kvs ...interface{}) {
	logw(criticalsv, l, msg, kvs)
}

// Alert logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Alert(v ...interface{}) {
	logm(alertsv, l, v...)
}

// Alertln logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Alertln(v ...interface{}) {
	logn(alertsv, l, v...)
}

// Alertf logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Alertf(format string, v ...interface{}) {
	logf(alertsv, l, format, v...)
}

// Alertj logs when a person must take an action immediately.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Alertj(msg string, v interface{}) {
	logj(alertsv, l, msg, v)
}

// Alertw logs when a person must take an action immediately.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Alertw(msg string, kvs ...interface{}) {
	logw(alertsv, l, msg, kvs)
}

// Emergency logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Emergency(v ...interface{}) {
	logm(emergencysv, l, v...)
}

// Emergencyln logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Emergencyln(v ...interface{}) {
	logn(emergencysv, l, v...)
}

// Emergencyf logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Emergencyf(format string, v ...interface{}) {
	logf(emergencysv, l, format, v...)
}

// Emergencyj logs when one or more systems are unusable.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Emergencyj(msg string, v interface{}) {
	logj(emergencysv, l, msg, v)
}

// Emergencyw logs when one or more systems are unusable.
// Arguments populate jsonPayload in the log entry.
func (l Logger) Emergencyw(msg string, kvs ...interface{}) {
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

func logm(s severity, l Logger, v ...interface{}) {
	logs(s, l, fmt.Sprint(v...))
}

func logn(s severity, l Logger, v ...interface{}) {
	logs(s, l, fmt.Sprintln(v...))
}

func logf(s severity, l Logger, format string, v ...interface{}) {
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

func logj(s severity, l Logger, msg string, j interface{}) {
	entry := make(map[string]json.RawMessage)
	if buf, err := json.Marshal(j); err != nil {
		panic(err)
	} else if err := json.Unmarshal(buf, &entry); err != nil {
		panic(err)
	}

	loge(s, l, msg, entry)
}

func logw(s severity, l Logger, msg string, kvs []interface{}) {
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
