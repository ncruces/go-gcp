// Package glog implements basic logging for Google Cloud Run
// and Cloud Functions.
package glog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var ProjectID string = os.Getenv("GOOGLE_CLOUD_PROJECT")

// Print logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	log(defaultsv, Logger{}, v...)
}

// Println logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	logln(defaultsv, Logger{}, v...)
}

// Printf logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	logf(defaultsv, Logger{}, format, v...)
}

// Debug logs debug or trace information.
// Arguments are handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	log(debugsv, Logger{}, v...)
}

// Debugln logs debug or trace information.
// Arguments are handled in the manner of fmt.Println.
func Debugln(v ...interface{}) {
	logln(debugsv, Logger{}, v...)
}

// Debugf logs debug or trace information.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	logf(debugsv, Logger{}, format, v...)
}

// Info logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	log(infosv, Logger{}, v...)
}

// Infoln logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Println.
func Infoln(v ...interface{}) {
	logln(infosv, Logger{}, v...)
}

// Infof logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	logf(infosv, Logger{}, format, v...)
}

// Notice logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Print.
func Notice(v ...interface{}) {
	log(noticesv, Logger{}, v...)
}

// Noticeln logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Println.
func Noticeln(v ...interface{}) {
	logln(noticesv, Logger{}, v...)
}

// Noticef logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Printf.
func Noticef(format string, v ...interface{}) {
	logf(noticesv, Logger{}, format, v...)
}

// Warning logs events that might cause problems.
// Arguments are handled in the manner of fmt.Print.
func Warning(v ...interface{}) {
	log(warningsv, Logger{}, v...)
}

// Warningln logs events that might cause problems.
// Arguments are handled in the manner of fmt.Println.
func Warningln(v ...interface{}) {
	logln(warningsv, Logger{}, v...)
}

// Warningf logs events that might cause problems.
// Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, v ...interface{}) {
	logf(warningsv, Logger{}, format, v...)
}

// Error logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	log(errorsv, Logger{}, v...)
}

// Errorln logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	logln(errorsv, Logger{}, v...)
}

// Errorf logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	logf(errorsv, Logger{}, format, v...)
}

// Critical logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Print.
func Critical(v ...interface{}) {
	log(criticalsv, Logger{}, v...)
}

// Criticalln logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Println.
func Criticalln(v ...interface{}) {
	logln(criticalsv, Logger{}, v...)
}

// Criticalf logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Printf.
func Criticalf(format string, v ...interface{}) {
	logf(criticalsv, Logger{}, format, v...)
}

// Alert logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Print.
func Alert(v ...interface{}) {
	log(alertsv, Logger{}, v...)
}

// Alertln logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Println.
func Alertln(v ...interface{}) {
	logln(alertsv, Logger{}, v...)
}

// Alertf logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Printf.
func Alertf(format string, v ...interface{}) {
	logf(alertsv, Logger{}, format, v...)
}

// Emergency logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Print.
func Emergency(v ...interface{}) {
	log(emergencysv, Logger{}, v...)
}

// Emergencyln logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Println.
func Emergencyln(v ...interface{}) {
	logln(emergencysv, Logger{}, v...)
}

// Emergencyf logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Printf.
func Emergencyf(format string, v ...interface{}) {
	logf(emergencysv, Logger{}, format, v...)
}

type Logger struct {
	trace  string
	spanID string
}

func ForRequest(r *http.Request) (l Logger) {
	if ProjectID != "" {
		h := r.Header.Get("X-Cloud-Trace-Context")

		i := strings.IndexByte(h, '/')
		if i > 0 {
			l.trace = fmt.Sprintf("projects/%s/traces/%s", ProjectID, h[0:i])
			j := strings.IndexByte(h[i:], ';')
			if j > 0 {
				l.spanID = h[i : i+j]
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

func log(s severity, l Logger, v ...interface{}) {
	logs(s, l, fmt.Sprint(v...))
}

func logln(s severity, l Logger, v ...interface{}) {
	logs(s, l, fmt.Sprintln(v...))
}

func logf(s severity, l Logger, format string, v ...interface{}) {
	logs(s, l, fmt.Sprintf(format, v...))
}

type entry struct {
	Message  string `json:"message"`
	Severity string `json:"severity,omitempty"`
	Trace    string `json:"logging.googleapis.com/trace,omitempty"`
	SpanID   string `json:"logging.googleapis.com/spanId,omitempty"`
}

func logs(s severity, l Logger, msg string) {
	var f *os.File
	if s >= errorsv {
		f = os.Stderr
	} else {
		f = os.Stdout
	}
	json.NewEncoder(f).Encode(entry{msg, s.String(), l.trace, l.spanID})
}
