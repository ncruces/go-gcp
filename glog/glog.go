// Package glog implements structured logging for Google App Engine, Cloud Run
// and Cloud Functions.
package glog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var std Logger

// ProjectID should be set to the Google Cloud project ID.
var ProjectID string = os.Getenv("GOOGLE_CLOUD_PROJECT")

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
// Arguments become jsonPayload in the log entry.
func Printj(msg string, v interface{}) {
	std.Printj(msg, v)
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
// Arguments become jsonPayload in the log entry.
func Debugj(msg string, v interface{}) {
	std.Debugj(msg, v)
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
// Arguments become jsonPayload in the log entry.
func Infoj(msg string, v interface{}) {
	std.Infoj(msg, v)
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
// Arguments become jsonPayload in the log entry.
func Noticej(msg string, v interface{}) {
	std.Noticej(msg, v)
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
// Arguments become jsonPayload in the log entry.
func Warningj(msg string, v interface{}) {
	std.Warningj(msg, v)
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
// Arguments become jsonPayload in the log entry.
func Errorj(msg string, v interface{}) {
	std.Errorj(msg, v)
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
// Arguments become jsonPayload in the log entry.
func Criticalj(msg string, v interface{}) {
	std.Criticalj(msg, v)
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
// Arguments become jsonPayload in the log entry.
func Alertj(msg string, v interface{}) {
	std.Alertj(msg, v)
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
// Arguments become jsonPayload in the log entry.
func Emergencyj(msg string, v interface{}) {
	std.Emergencyj(msg, v)
}

type Logger struct {
	trace  string
	spanID string
}

func ForRequest(r *http.Request) (l Logger) {
	if ProjectID != "" {
		h := r.Header.Get("X-Cloud-Trace-Context")
		if i := strings.IndexByte(h, '/'); i > 0 {
			if t := h[:i]; strings.Count(t, "0") != len(t) {
				l.trace = fmt.Sprintf("projects/%s/traces/%s", ProjectID, t)
			}
		}
	}
	return l
}

// Print logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Print(v ...interface{}) {
	log(defaultsv, l, v...)
}

// Println logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Println(v ...interface{}) {
	logln(defaultsv, l, v...)
}

// Printf logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Printf(format string, v ...interface{}) {
	logf(defaultsv, l, format, v...)
}

// Printj logs an entry with no assigned severity level.
// Arguments become jsonPayload in the log entry.
func (l Logger) Printj(msg string, v interface{}) {
	logj(defaultsv, l, msg, v)
}

// Debug logs debug or trace information.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Debug(v ...interface{}) {
	log(debugsv, l, v...)
}

// Debugln logs debug or trace information.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Debugln(v ...interface{}) {
	logln(debugsv, l, v...)
}

// Debugf logs debug or trace information.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Debugf(format string, v ...interface{}) {
	logf(debugsv, l, format, v...)
}

// Debugj logs debug or trace information.
// Arguments become jsonPayload in the log entry.
func (l Logger) Debugj(msg string, v interface{}) {
	logj(debugsv, l, msg, v)
}

// Info logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Info(v ...interface{}) {
	log(infosv, l, v...)
}

// Infoln logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Infoln(v ...interface{}) {
	logln(infosv, l, v...)
}

// Infof logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Infof(format string, v ...interface{}) {
	logf(infosv, l, format, v...)
}

// Infoj logs routine information, such as ongoing status or performance.
// Arguments become jsonPayload in the log entry.
func (l Logger) Infoj(msg string, v interface{}) {
	logj(infosv, l, msg, v)
}

// Notice logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Notice(v ...interface{}) {
	log(noticesv, l, v...)
}

// Noticeln logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Noticeln(v ...interface{}) {
	logln(noticesv, l, v...)
}

// Noticef logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Noticef(format string, v ...interface{}) {
	logf(noticesv, l, format, v...)
}

// Noticej logs normal but significant events, such as start up, shut down, or configuration.
// Arguments become jsonPayload in the log entry.
func (l Logger) Noticej(msg string, v interface{}) {
	logj(noticesv, l, msg, v)
}

// Warning logs events that might cause problems.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Warning(v ...interface{}) {
	log(warningsv, l, v...)
}

// Warningln logs events that might cause problems.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Warningln(v ...interface{}) {
	logln(warningsv, l, v...)
}

// Warningf logs events that might cause problems.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Warningf(format string, v ...interface{}) {
	logf(warningsv, l, format, v...)
}

// Warningj logs events that might cause problems.
// Arguments become jsonPayload in the log entry.
func (l Logger) Warningj(msg string, v interface{}) {
	logj(warningsv, l, msg, v)
}

// Error logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Error(v ...interface{}) {
	log(errorsv, l, v...)
}

// Errorln logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Errorln(v ...interface{}) {
	logln(errorsv, l, v...)
}

// Errorf logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Errorf(format string, v ...interface{}) {
	logf(errorsv, l, format, v...)
}

// Errorj logs events likely to cause problems.
// Arguments become jsonPayload in the log entry.
func (l Logger) Errorj(msg string, v interface{}) {
	logj(errorsv, l, msg, v)
}

// Critical logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Critical(v ...interface{}) {
	log(criticalsv, l, v...)
}

// Criticalln logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Criticalln(v ...interface{}) {
	logln(criticalsv, l, v...)
}

// Criticalf logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Criticalf(format string, v ...interface{}) {
	logf(criticalsv, l, format, v...)
}

// Criticalj logs events that cause more severe problems or outages.
// Arguments become jsonPayload in the log entry.
func (l Logger) Criticalj(msg string, v interface{}) {
	logj(criticalsv, l, msg, v)
}

// Alert logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Alert(v ...interface{}) {
	log(alertsv, l, v...)
}

// Alertln logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Alertln(v ...interface{}) {
	logln(alertsv, l, v...)
}

// Alertf logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Alertf(format string, v ...interface{}) {
	logf(alertsv, l, format, v...)
}

// Alertj logs when a person must take an action immediately.
// Arguments become jsonPayload in the log entry.
func (l Logger) Alertj(msg string, v interface{}) {
	logj(alertsv, l, msg, v)
}

// Emergency logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Print.
func (l Logger) Emergency(v ...interface{}) {
	log(emergencysv, l, v...)
}

// Emergencyln logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Println.
func (l Logger) Emergencyln(v ...interface{}) {
	logln(emergencysv, l, v...)
}

// Emergencyf logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Printf.
func (l Logger) Emergencyf(format string, v ...interface{}) {
	logf(emergencysv, l, format, v...)
}

// Emergencyj logs when one or more systems are unusable.
// Arguments become jsonPayload in the log entry.
func (l Logger) Emergencyj(msg string, v interface{}) {
	logj(emergencysv, l, msg, v)
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
	}
}

func (s severity) File() *os.File {
	if s >= errorsv {
		return os.Stderr
	} else {
		return os.Stdout
	}
}

func log(s severity, l Logger, v ...interface{}) {
	logs(s, l, fmt.Sprint(v...))
}

func logln(s severity, l Logger, v ...interface{}) {
	logs(s, l, fmt.Sprintln(v...))
}

func logf(s severity, l Logger, format string, v ...interface{}) {
	logs(s, l, fmt.Sprintf(format, v...))
}

func logs(s severity, l Logger, msg string) {
	entry := entry{msg, s.String(), l.trace}
	json.NewEncoder(s.File()).Encode(entry)
}

func logj(s severity, l Logger, msg string, j interface{}) {
	entry := make(map[string]json.RawMessage)
	if buf, err := json.Marshal(j); err != nil {
		panic(err)
	} else if err := json.Unmarshal(buf, &entry); err != nil {
		panic(err)
	}

	if v := msg; v != "" {
		entry["message"], _ = json.Marshal(v)
	}
	if v := s.String(); v != "" {
		entry["severity"], _ = json.Marshal(v)
	}
	if v := l.trace; v != "" {
		entry["logging.googleapis.com/trace"], _ = json.Marshal(v)
	}

	json.NewEncoder(s.File()).Encode(entry)
}

type entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`
}
