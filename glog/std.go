package glog

import "log"

type stdLogger struct{}

func (s stdLogger) Write(p []byte) (int, error) {
	logs(defaultsv, std, string(p))
	return len(p), nil
}

// SetupLogger sets up a log.Logger to output structured logs at the default severity level.
func SetupLogger(l *log.Logger) {
	l.SetFlags(0)
	l.SetOutput(stdLogger{})
}
