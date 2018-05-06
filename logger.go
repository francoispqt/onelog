package onelog

import (
	"io"

	"github.com/francoispqt/gojay"
)

var msgKey = "message"

// LevelText personalises the text for a specific level
func LevelText(level int, txt string) {
	mux.Lock()
	Levels[level] = txt
	genLevelSlices()
	genStreamLevelSlices()
	mux.Unlock()
}

// MsgKey sets the key for the message field
func MsgKey(s string) {
	mux.Lock()
	msgKey = s
	genLevelSlices()
	genStreamLevelSlices()
}

// LevelKey sets the key for the level field
func LevelKey(s string) {
	mux.Lock()
	levelKey = s
	genLevelSlices()
	genStreamLevelSlices()
	mux.Unlock()
}

// Encoder is an alias to gojay.Encoder
type Encoder = gojay.Encoder

// ObjectFunc is an alias to gojay.EncodeObjectFunc
type ObjectFunc = gojay.EncodeObjectFunc

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

// logging functions

// // Log logs an entry
// func (l *Logger) Log(level uint8, msg string) {
// 	// first find writer for level
// 	// if none, stop
// 	if level&l.levels == 0 {
// 		return
// 	}

// 	// then call format on formatter
// 	l.Format(func(f *Encoder) {
// 		f.AppendStringKey(l.lvlKey, Levels[level])
// 		f.AppendString(msg)
// 		if l.hook != nil {
// 			l.hook(f)
// 		}
// 	})
// }

// // LogWithFields logs an entry with custom fields
// func (l *Logger) LogWithFields(level uint8, msg string, fields func(*Encoder)) {
// 	// first find writer for level
// 	// if none, stop
// 	if level&l.levels == 0 {
// 		return
// 	}

// 	// then call format on formatter
// 	l.Format(func(f *Encoder) {
// 		f.AppendStringKey(l.lvlKey, Levels[level])
// 		f.AppendString(msg)
// 		if l.hook != nil {
// 			l.hook(f)
// 		}
// 		fields(f)
// 	})
// }

// Info logs an entry with INFO level
func (l *Logger) Info(msg string) {
	// first find writer for level
	// if none, stop
	if INFO&l.levels == 0 {
		return
	}

	// then call format on formatter
	l.Format(INFO, msg)
}

// InfoWithFields logs an entry with INFO level and custom fields
func (l *Logger) InfoWithFields(msg string, fields func(*Encoder)) {
	// first find writer for level
	// if none, stop
	if INFO&l.levels == 0 {
		return
	}

	// then call format on formatter
	l.FormatFields(INFO, msg, fields)
}

// Debug logs an entry with DEBUG level
func (l *Logger) Debug(msg string) {
	// check if level is in config
	// if not, return
	if DEBUG&l.levels == 0 {
		return
	}
	// then call format on formatter
	l.Format(DEBUG, msg)
}

// DebugWithFields logs an entry with DEBUG level and custom fields
func (l *Logger) DebugWithFields(msg string, fields func(*Encoder)) {
	// check if level is in config
	// if not, return
	if DEBUG&l.levels == 0 {
		return
	}
	// then call format on formatter
	l.FormatFields(DEBUG, msg, fields)
}

// Warn logs an entry with WARN level
func (l *Logger) Warn(msg string) {
	// check if level is in config
	// if not, return
	if WARN&l.levels == 0 {
		return
	}
	// then call format on formatter
	l.Format(WARN, msg)
}

// WarnWithFields logs an entry with WARN level and custom fields
func (l *Logger) WarnWithFields(msg string, fields func(*Encoder)) {
	// check if level is in config
	// if not, return
	if WARN&l.levels == 0 {
		return
	}
	// then call format on formatter
	l.FormatFields(WARN, msg, fields)
}

// Error logs an entry with ERROR level
func (l *Logger) Error(msg string) {
	// check if level is in config
	// if not, return
	if ERROR&l.levels == 0 {
		return
	}
	// then call format on formatter
	l.Format(ERROR, msg)
}

// ErrorWithFields logs an entry with ERROR level and custom fields
func (l *Logger) ErrorWithFields(msg string, fields func(*Encoder)) {
	// check if level is in config
	// if not, return
	if ERROR&l.levels == 0 {
		return
	}
	// then call format on formatter
	l.FormatFields(ERROR, msg, fields)
}

// Fatal logs an entry with FATAL level
func (l *Logger) Fatal(msg string) {
	// check if level is in config
	// if not, return
	if FATAL&l.levels == 0 {
		return
	}
	// then call format on formatter
	l.Format(FATAL, msg)
}

// FatalWithFields logs an entry with FATAL level and custom fields
func (l *Logger) FatalWithFields(msg string, fields func(*Encoder)) {
	// check if level is in config
	// if not, return
	if FATAL&l.levels == 0 {
		return
	}
	// then call format on formatter
	l.FormatFields(FATAL, msg, fields)
}

func (l *Logger) Format(level int, msg string) {
	// get a decoder
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	// append first part containing object initialisation
	// and message key
	enc.AppendBuf(levelsJSON[level])
	enc.AppendString(msg)
	if l.ctx != nil {
		enc.AppendBuf(l.ctx)
	}
	if l.hook != nil {
		l.hook(enc)
	}
	enc.AppendByte('}')
	enc.Write()
}

// FormatFields formats the log entry by calling Encoder.EncodeObject()
func (l *Logger) FormatFields(level int, msg string, fields func(*Encoder)) {
	// get a decoder
	enc := gojay.BorrowEncoder(l.w)
	defer enc.Release()
	// append first part containing object initialisation
	// and message key
	enc.AppendBuf(levelsJSON[level])
	enc.AppendString(msg)
	if l.ctx != nil {
		enc.AppendBuf(l.ctx)
	}
	if l.hook != nil {
		l.hook(enc)
	}
	fields(enc)
	enc.AppendByte('}')
	enc.Write()
}
