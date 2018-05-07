package log

import (
	"os"
	"time"

	"github.com/francoispqt/onelog"
)

var logger = onelog.New(os.Stdout, onelog.ALL).Hook(func(e onelog.Entry) {
	e.Int64("time", time.Now().Unix())
})

// Info prints a message with log level Info.
func Info(msg string) {
	logger.Info(msg)
}

// InfoWithFields prints a message with log level INFO and fields.
func InfoWithFields(msg string, fields func(e onelog.Entry)) {
	logger.InfoWithFields(msg, fields)
}

// Debug prints a message with log level DEBUG.
func Debug(msg string) {
	logger.Debug(msg)
}

// DebugWithFields prints a message with log level DEBUG and fields.
func DebugWithFields(msg string, fields func(e onelog.Entry)) {
	logger.DebugWithFields(msg, fields)
}

// Warn prints a message with log level INFO.
func Warn(msg string) {
	logger.Warn(msg)
}

// WarnWithFields prints a message with log level WARN and fields.
func WarnWithFields(msg string, fields func(e onelog.Entry)) {
	logger.WarnWithFields(msg, fields)
}

// Error prints a message with log level ERROR.
func Error(msg string) {
	logger.Error(msg)
}

// ErrorWithFields prints a message with log level ERROR and fields.
func ErrorWithFields(msg string, fields func(e onelog.Entry)) {
	logger.ErrorWithFields(msg, fields)
}

// Fatal prints a message with log level FATAL.
func Fatal(msg string) {
	logger.Fatal(msg)
}

// FatalWithFields prints a message with log level FATAL and fields.
func FatalWithFields(msg string, fields func(e onelog.Entry)) {
	logger.FatalWithFields(msg, fields)
}
