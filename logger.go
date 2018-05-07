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
	hook   func(Entry)
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
func (l *Logger) Hook(h func(Entry)) *Logger {
	l.hook = h
	return l
}

func (l *Logger) copy() *Logger {
	nL := Logger{
		levels: l.levels,
		w:      l.w,
	}
	return &nL
}

// With copys the current Logger and adds it a context
func (l *Logger) With(f func(Entry)) *Logger {
	nL := l.copy()
	e := Entry{}
	enc := gojay.BorrowEncoder(nL.w)
	e.enc = enc
	enc.AppendByte(' ')
	f(e)
	nL.ctx = enc.Buf()[1:]
	enc.Release()
	return nL
}

// Info logs an entry with INFO level
func (l *Logger) Info(msg string) {
	// first find writer for level
	// if none, stop
	if INFO&l.levels == 0 {
		return
	}
	e := Entry{}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc
	l.beginLogEntry(INFO, msg, enc)
	if l.hook != nil {
		l.hook(e)
	}
	l.closeLogEntry(enc)
	enc.Release()
}

// InfoWithFields logs an entry with INFO level and custom fields
func (l *Logger) InfoWithFields(msg string, fields func(Entry)) {
	// first find writer for level
	// if none, stop
	if INFO&l.levels == 0 {
		return
	}
	e := Entry{}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc
	l.beginLogEntry(INFO, msg, enc)
	if l.hook != nil {
		l.hook(e)
	}
	fields(e)
	l.closeLogEntry(enc)
	enc.Release()
}

// Debug logs an entry with DEBUG level
func (l *Logger) Debug(msg string) {
	// check if level is in config
	// if not, return
	if DEBUG&l.levels == 0 {
		return
	}
	e := Entry{}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc
	l.beginLogEntry(DEBUG, msg, enc)
	if l.hook != nil {
		l.hook(e)
	}
	l.closeLogEntry(enc)
	enc.Release()
}

// DebugWithFields logs an entry with DEBUG level and custom fields
func (l *Logger) DebugWithFields(msg string, fields func(Entry)) {
	// check if level is in config
	// if not, return
	if DEBUG&l.levels == 0 {
		return
	}
	e := Entry{}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc
	l.beginLogEntry(DEBUG, msg, enc)
	if l.hook != nil {
		l.hook(e)
	}
	fields(e)
	l.closeLogEntry(enc)
	enc.Release()
}

// Warn logs an entry with WARN level
func (l *Logger) Warn(msg string) {
	// check if level is in config
	// if not, return
	if WARN&l.levels == 0 {
		return
	}
	e := Entry{}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc
	l.beginLogEntry(WARN, msg, enc)
	if l.hook != nil {
		l.hook(e)
	}
	l.closeLogEntry(enc)
	enc.Release()
}

// WarnWithFields logs an entry with WARN level and custom fields
func (l *Logger) WarnWithFields(msg string, fields func(Entry)) {
	if WARN&l.levels == 0 {
		return
	}
	e := Entry{}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc
	l.beginLogEntry(WARN, msg, enc)
	if l.hook != nil {
		l.hook(e)
	}
	fields(e)
	l.closeLogEntry(enc)
	enc.Release()
}

// Error logs an entry with ERROR level
func (l *Logger) Error(msg string) {
	if ERROR&l.levels == 0 {
		return
	}
	e := Entry{}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc
	l.beginLogEntry(ERROR, msg, enc)
	if l.hook != nil {
		l.hook(e)
	}
	l.closeLogEntry(enc)
	enc.Release()
}

// ErrorWithFields logs an entry with ERROR level and custom fields
func (l *Logger) ErrorWithFields(msg string, fields func(Entry)) {
	if ERROR&l.levels == 0 {
		return
	}
	e := Entry{}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc
	l.beginLogEntry(ERROR, msg, enc)
	if l.hook != nil {
		l.hook(e)
	}
	fields(e)
	l.closeLogEntry(enc)
	enc.Release()
}

// Fatal logs an entry with FATAL level
func (l *Logger) Fatal(msg string) {
	if FATAL&l.levels == 0 {
		return
	}
	e := Entry{}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc
	l.beginLogEntry(FATAL, msg, enc)
	if l.hook != nil {
		l.hook(e)
	}
	l.closeLogEntry(enc)
	enc.Release()
}

// FatalWithFields logs an entry with FATAL level and custom fields
func (l *Logger) FatalWithFields(msg string, fields func(Entry)) {
	if FATAL&l.levels == 0 {
		return
	}
	e := Entry{}
	// then call format on formatter
	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc
	l.beginLogEntry(FATAL, msg, enc)
	if l.hook != nil {
		l.hook(e)
	}
	fields(e)
	l.closeLogEntry(enc)
	enc.Release()
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
