package onelog

import (
	"io"
	"runtime"
	"strconv"

	"github.com/francoispqt/gojay"
)

var logClose = []byte("}\n")
var msgKey = "message"

// LevelText personalises the text for a specific level
func LevelText(level int, txt string) {
	Levels[level] = txt
	genLevelSlices()
}

// MsgKey sets the key for the message field
func MsgKey(s string) {
	msgKey = s
	genLevelSlices()
}

// LevelKey sets the key for the level field
func LevelKey(s string) {
	levelKey = s
	genLevelSlices()
}

// Encoder is an alias to gojay.Encoder
type Encoder = gojay.Encoder

// Object is an alias to gojay.EncodeObjectFunc
type Object = gojay.EncodeObjectFunc

// Logger is the type representing a logger.
type Logger struct {
	hook   func(*Encoder)
	w      io.Writer
	levels uint8
	ctx    []byte
}

// NewLogger returns a fresh onelog Logger with default values.
func NewLogger(w io.Writer, levels uint8) *Logger {
	return &Logger{
		levels: levels,
		w:      w,
	}
}

// Hook sets a hook to run for all log entries to add generic fields
func (l *Logger) Hook(h func(*Encoder)) {
	l.hook = h
}

func (l *Logger) copy() *Logger {
	nL := Logger{
		levels: l.levels,
		w:      l.w,
	}
	return &nL
}

// With copys the current Logger and adds it a context
func (l *Logger) With(f func(*Encoder)) *Logger {
	nL := l.copy()
	enc := gojay.BorrowEncoder(nL.w)
	defer enc.Release()
	enc.AppendByte(' ')
	f(enc)
	nL.ctx = enc.Buf()[1:]
	return nL
}

// Info logs an entry with INFO level
func (l *Logger) Info(msg string) {
	// first find writer for level
	// if none, stop
	if INFO&l.levels == 0 {
		return
	}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	l.beginLogEntry(INFO, msg, enc)
	if l.hook != nil {
		l.hook(enc)
	}
	l.closeLogEntry(enc)
}

// InfoWithFields logs an entry with INFO level and custom fields
func (l *Logger) InfoWithFields(msg string, fields func(*Encoder)) {
	// first find writer for level
	// if none, stop
	if INFO&l.levels == 0 {
		return
	}

	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	l.beginLogEntry(INFO, msg, enc)
	if l.hook != nil {
		l.hook(enc)
	}
	fields(enc)
	l.closeLogEntry(enc)
}

// Debug logs an entry with DEBUG level
func (l *Logger) Debug(msg string) {
	// check if level is in config
	// if not, return
	if DEBUG&l.levels == 0 {
		return
	}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	l.beginLogEntry(DEBUG, msg, enc)
	if l.hook != nil {
		l.hook(enc)
	}
	l.closeLogEntry(enc)
}

// DebugWithFields logs an entry with DEBUG level and custom fields
func (l *Logger) DebugWithFields(msg string, fields func(*Encoder)) {
	// check if level is in config
	// if not, return
	if DEBUG&l.levels == 0 {
		return
	}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	l.beginLogEntry(DEBUG, msg, enc)
	if l.hook != nil {
		l.hook(enc)
	}
	fields(enc)
	l.closeLogEntry(enc)
}

// Warn logs an entry with WARN level
func (l *Logger) Warn(msg string) {
	// check if level is in config
	// if not, return
	if WARN&l.levels == 0 {
		return
	}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	l.beginLogEntry(WARN, msg, enc)
	if l.hook != nil {
		l.hook(enc)
	}
	l.closeLogEntry(enc)
}

// WarnWithFields logs an entry with WARN level and custom fields
func (l *Logger) WarnWithFields(msg string, fields func(*Encoder)) {
	// check if level is in config
	// if not, return
	if WARN&l.levels == 0 {
		return
	}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	l.beginLogEntry(WARN, msg, enc)
	if l.hook != nil {
		l.hook(enc)
	}
	fields(enc)
	l.closeLogEntry(enc)
}

// Error logs an entry with ERROR level
func (l *Logger) Error(msg string) {
	// check if level is in config
	// if not, return
	if ERROR&l.levels == 0 {
		return
	}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	l.beginLogEntry(ERROR, msg, enc)
	if l.hook != nil {
		l.hook(enc)
	}
	l.closeLogEntry(enc)
}

// ErrorWithFields logs an entry with ERROR level and custom fields
func (l *Logger) ErrorWithFields(msg string, fields func(*Encoder)) {
	// check if level is in config
	// if not, return
	if ERROR&l.levels == 0 {
		return
	}
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	l.beginLogEntry(ERROR, msg, enc)
	if l.hook != nil {
		l.hook(enc)
	}
	fields(enc)
	l.closeLogEntry(enc)
}

// Fatal logs an entry with FATAL level
func (l *Logger) Fatal(msg string) {
	// check if level is in config
	// if not, return
	if FATAL&l.levels == 0 {
		return
	}
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	l.beginLogEntry(FATAL, msg, enc)
	if l.hook != nil {
		l.hook(enc)
	}
	l.closeLogEntry(enc)
}

// FatalWithFields logs an entry with FATAL level and custom fields
func (l *Logger) FatalWithFields(msg string, fields func(*Encoder)) {
	// check if level is in config
	// if not, return
	if FATAL&l.levels == 0 {
		return
	}
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	l.beginLogEntry(FATAL, msg, enc)
	if l.hook != nil {
		l.hook(enc)
	}
	fields(enc)
	l.closeLogEntry(enc)
}

func (l *Logger) beginLogEntry(level int, msg string, enc *Encoder) {
	enc.AppendBytes(levelsJSON[level])
	enc.AppendString(msg)
	if l.ctx != nil {
		enc.AppendBytes(l.ctx)
	}
}

func (l Logger) closeLogEntry(enc *Encoder) {
	enc.AppendBytes(logClose)
	enc.Write()
}

// Caller returns the caller in the stack trace, skipped n times.
func (l *Logger) Caller(n int) string {
	_, f, fl, _ := runtime.Caller(n)
	flStr := strconv.Itoa(fl)
	return f + ":" + flStr
}
