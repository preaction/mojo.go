package mojo

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Log is a basic logger for applications
type Log struct {
	Handle  io.Writer
	Short   bool
	context interface{}
	level   int
}

var syslogPri = map[string]int{
	"debug": 7,
	"info":  6,
	"warn":  4,
	"error": 3,
	"fatal": 2,
}

// NewLog builds a log object with the default settings.
func NewLog() Log {
	level := os.Getenv("MOJO_LOG_LEVEL")
	if level == "" {
		level = "debug"
	}
	return Log{level: syslogPri[level]}
}

// Level allows changing the log level
func (log *Log) Level(level string) {
	log.level = syslogPri[level]
}

// Context creates a logger, appending the given context data to each
// log message
func (log *Log) Context(ctx interface{}) Log {
	return Log{context: ctx, Handle: log, level: log.level, Short: log.Short}
}

// Log a debug level message
func (log *Log) Debug(msg string, args ...interface{}) {
	if log.level < syslogPri["debug"] {
		return
	}
	log.message("debug", msg, args...)
}

// Log an info level message
func (log *Log) Info(msg string, args ...interface{}) {
	if log.level < syslogPri["info"] {
		return
	}
	log.message("info", msg, args...)
}

// Log a warn level message
func (log *Log) Warn(msg string, args ...interface{}) {
	if log.level < syslogPri["warn"] {
		return
	}
	log.message("warn", msg, args...)
}

// Log an error level message
func (log *Log) Error(msg string, args ...interface{}) {
	if log.level < syslogPri["error"] {
		return
	}
	log.message("error", msg, args...)
}

// Log a fatal level message
func (log *Log) Fatal(msg string, args ...interface{}) {
	if log.level < syslogPri["fatal"] {
		return
	}
	log.message("fatal", msg, args...)
}

// message handles writing to the handle, preparing any necessary
// context
func (log *Log) message(level string, msg string, args ...interface{}) {
	var format string
	if log.Short {
		format = fmt.Sprintf("<%d>[%s] ", syslogPri[level], level[0:1])
	} else {
		now := time.Now().Truncate(time.Millisecond).Format("2006-01-02 15:04:05.000")
		format = fmt.Sprintf("[%s] [%s] ", now, level)
	}

	if log.context != nil {
		format += fmt.Sprintf("%s ", log.context)
	}

	format += msg + "\n"
	log.Write([]byte(fmt.Sprintf(format, args...)))
}

// Write implements the Writer interface to allow for child Log objects
func (log *Log) Write(msg []byte) (int, error) {
	return log.Handle.Write(msg)
}
