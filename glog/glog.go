// Package glog implements basic logging for Google Cloud Run
// and Cloud Functions.
package glog

import (
	"encoding/json"
	"fmt"
	"os"
)

// Print logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	log(defaultsv, v...)
}

// Println logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	logln(defaultsv, v...)
}

// Printf logs an entry with no assigned severity level.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	logf(defaultsv, format, v...)
}

// Debug logs debug or trace information.
// Arguments are handled in the manner of fmt.Print.
func Debug(v ...interface{}) {
	log(debugsv, v...)
}

// Debugln logs debug or trace information.
// Arguments are handled in the manner of fmt.Println.
func Debugln(v ...interface{}) {
	logln(debugsv, v...)
}

// Debugf logs debug or trace information.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	logf(debugsv, format, v...)
}

// Info logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Print.
func Info(v ...interface{}) {
	log(infosv, v...)
}

// Infoln logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Println.
func Infoln(v ...interface{}) {
	logln(infosv, v...)
}

// Infof logs routine information, such as ongoing status or performance.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	logf(infosv, format, v...)
}

// Notice logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Print.
func Notice(v ...interface{}) {
	log(noticesv, v...)
}

// Noticeln logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Println.
func Noticeln(v ...interface{}) {
	logln(noticesv, v...)
}

// Noticef logs normal but significant events, such as start up, shut down, or configuration.
// Arguments are handled in the manner of fmt.Printf.
func Noticef(format string, v ...interface{}) {
	logf(noticesv, format, v...)
}

// Warning logs events that might cause problems.
// Arguments are handled in the manner of fmt.Print.
func Warning(v ...interface{}) {
	log(warningsv, v...)
}

// Warningln logs events that might cause problems.
// Arguments are handled in the manner of fmt.Println.
func Warningln(v ...interface{}) {
	logln(warningsv, v...)
}

// Warningf logs events that might cause problems.
// Arguments are handled in the manner of fmt.Printf.
func Warningf(format string, v ...interface{}) {
	logf(warningsv, format, v...)
}

// Error logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	log(errorsv, v...)
}

// Errorln logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	logln(errorsv, v...)
}

// Errorf logs events likely to cause problems.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	logf(errorsv, format, v...)
}

// Critical logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Print.
func Critical(v ...interface{}) {
	log(criticalsv, v...)
}

// Criticalln logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Println.
func Criticalln(v ...interface{}) {
	logln(criticalsv, v...)
}

// Criticalf logs events that cause more severe problems or outages.
// Arguments are handled in the manner of fmt.Printf.
func Criticalf(format string, v ...interface{}) {
	logf(criticalsv, format, v...)
}

// Alert logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Print.
func Alert(v ...interface{}) {
	log(alertsv, v...)
}

// Alertln logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Println.
func Alertln(v ...interface{}) {
	logln(alertsv, v...)
}

// Alertf logs when a person must take an action immediately.
// Arguments are handled in the manner of fmt.Printf.
func Alertf(format string, v ...interface{}) {
	logf(alertsv, format, v...)

}

// Emergency logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Print.
func Emergency(v ...interface{}) {
	log(emergencysv, v...)
}

// Emergencyln logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Println.
func Emergencyln(v ...interface{}) {
	logln(emergencysv, v...)
}

// Emergencyf logs when one or more systems are unusable.
// Arguments are handled in the manner of fmt.Printf.
func Emergencyf(format string, v ...interface{}) {
	logf(emergencysv, format, v...)
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

func (s severity) MarshalText() (text []byte, err error) {
	switch s {
	default:
		return nil, nil
	case debugsv:
		return []byte("DEBUG"), nil
	case infosv:
		return []byte("INFO"), nil
	case noticesv:
		return []byte("NOTICE"), nil
	case warningsv:
		return []byte("WARNING"), nil
	case errorsv:
		return []byte("ERROR"), nil
	case criticalsv:
		return []byte("CRITICAL"), nil
	case alertsv:
		return []byte("ALERT"), nil
	}
}

func log(s severity, v ...interface{}) {
	logs(s, fmt.Sprint(v...))
}

func logln(s severity, v ...interface{}) {
	logs(s, fmt.Sprintln(v...))
}

func logf(s severity, format string, v ...interface{}) {
	logs(s, fmt.Sprintf(format, v...))
}

type entry struct {
	Message  string   `json:"message"`
	Severity severity `json:"severity,omitempty"`
}

func logs(s severity, msg string) {
	var f *os.File
	if s >= errorsv {
		f = os.Stderr
	} else {
		f = os.Stdout
	}
	json.NewEncoder(f).Encode(entry{msg, s})
}
