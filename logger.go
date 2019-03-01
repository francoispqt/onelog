package onelog

import (
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"

	"github.com/francoispqt/gojay"
)

var (
	logOpen      = []byte("{")
	logClose     = []byte("}\n")
	logCloseOnly = []byte("}")
	msgKey       = "message"
)

// LevelText personalises the text for a specific level.
func LevelText(level uint8, txt string) {
	Levels[level] = txt
	genLevelSlices()
}

// MsgKey sets the key for the message field.
func MsgKey(s string) {
	msgKey = s
	genLevelSlices()
}

// LevelKey sets the key for the level field.
func LevelKey(s string) {
	levelKey = s
	genLevelSlices()
}

// Encoder is an alias to gojay.Encoder.
type Encoder = gojay.Encoder

// Object is an alias to gojay.EncodeObjectFunc.
type Object = gojay.EncodeObjectFunc

// ExitFunc is used to exit the app, `os.Exit()` is set as default on `New()`
type ExitFunc func(int)

// Logger is the type representing a logger.
type Logger struct {
	hook        func(Entry)
	w           io.Writer
	levels      uint8
	ctx         []func(Entry)
	ExitFn      ExitFunc
	contextName string
}

// New returns a fresh onelog Logger with default values.
func New(w io.Writer, levels uint8) *Logger {
	if w == nil {
		w = ioutil.Discard
	}

	return &Logger{
		w:      w,
		levels: levels,
		ExitFn: os.Exit,
	}
}

// NewContext returns a fresh onelog Logger with default values and
// context name set to provided contextName value.
func NewContext(w io.Writer, levels uint8, contextName string) *Logger {
	if w == nil {
		w = ioutil.Discard
	}

	return &Logger{
		w:           w,
		levels:      levels,
		contextName: contextName,
		ExitFn:      os.Exit,
	}
}

// Hook sets a hook to run for all log entries to add generic fields
func (l *Logger) Hook(h func(Entry)) *Logger {
	l.hook = h
	return l
}

func (l *Logger) copy(ctxName string) *Logger {
	nL := &Logger{
		levels:      l.levels,
		w:           l.w,
		hook:        l.hook,
		contextName: ctxName,
		ExitFn:      l.ExitFn,
	}
	if len(l.ctx) > 0 {
		var ctx = make([]func(e Entry), len(l.ctx))
		copy(ctx, l.ctx)
		nL.ctx = ctx
	}
	return nL
}

// With copies the current Logger and adds it a given context by running func f.
func (l *Logger) With(f func(Entry)) *Logger {
	nL := l.copy(l.contextName)

	if len(nL.ctx) == 0 {
		nL.ctx = make([]func(Entry), 0, 1)
	}

	nL.ctx = append(nL.ctx, f)
	return nL
}

// WithContext copies current logger enforcing all entry fields to be
// set into a map with the contextName set as the key name for giving map.
// This allows allocating all future uses of the logging methods to
// follow such formatting. The only exception are values provided by
// added hooks which will remain within the root level of generated json.
func (l *Logger) WithContext(contextName string) *Logger {
	nl := l.copy(contextName)
	return nl
}

// Info logs an entry with INFO level.
func (l *Logger) Info(msg string) {
	// first find writer for level
	// if none, stop
	if INFO&l.levels == 0 {
		return
	}
	e := Entry{Level: INFO, Message: msg}

	enc := gojay.BorrowEncoder(l.w)
	e.enc = enc

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e)
		l.runHook(e)
	} else {
		l.openEntry(enc)
	}

	l.closeEntry(e)
	l.finalizeIfContext(e)

	enc.Release()
}

// InfoWith return an ChainEntry with INFO level.
func (l *Logger) InfoWith(msg string) ChainEntry {
	// first find writer for level
	// if none, stop
	e := ChainEntry{
		Entry: Entry{
			l:       l,
			Level:   INFO,
			Message: msg,
		},
	}
	e.disabled = INFO&e.l.levels == 0
	if e.disabled {
		return e
	}

	e.Entry.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e.Entry)
		l.runHook(e.Entry)
		return e
	}

	l.openEntry(e.Entry.enc)
	return e
}

// InfoWithFields logs an entry with INFO level and custom fields.
func (l *Logger) InfoWithFields(msg string, fields func(Entry)) {
	// first find writer for level
	// if none, stop
	if INFO&l.levels == 0 {
		return
	}
	e := Entry{Level: INFO, Message: msg}

	e.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e)
		l.runHook(e)
	} else {
		l.openEntry(e.enc)
	}

	fields(e)
	l.closeEntry(e)
	l.finalizeIfContext(e)

	e.enc.Release()
}

// Debug logs an entry with DEBUG level.
func (l *Logger) Debug(msg string) {
	// check if level is in config
	// if not, return
	if DEBUG&l.levels == 0 {
		return
	}
	e := Entry{Level: DEBUG, Message: msg}

	e.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e)
		l.runHook(e)
	} else {
		l.openEntry(e.enc)
	}

	l.closeEntry(e)
	l.finalizeIfContext(e)

	e.enc.Release()
}

// DebugWith return ChainEntry with DEBUG level.
func (l *Logger) DebugWith(msg string) ChainEntry {
	// first find writer for level
	// if none, stop
	e := ChainEntry{
		Entry: Entry{
			l:       l,
			Level:   DEBUG,
			Message: msg,
		},
	}
	e.disabled = DEBUG&e.l.levels == 0
	if e.disabled {
		return e
	}

	e.Entry.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e.Entry)
		l.runHook(e.Entry)
		return e
	}

	l.openEntry(e.Entry.enc)
	return e
}

// DebugWithFields logs an entry with DEBUG level and custom fields.
func (l *Logger) DebugWithFields(msg string, fields func(Entry)) {
	// check if level is in config
	// if not, return
	if DEBUG&l.levels == 0 {
		return
	}
	e := Entry{Level: DEBUG, Message: msg}

	e.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e)
		l.runHook(e)
	} else {
		l.openEntry(e.enc)
	}

	fields(e)
	l.closeEntry(e)
	l.finalizeIfContext(e)

	e.enc.Release()
}

// Warn logs an entry with WARN level.
func (l *Logger) Warn(msg string) {
	// check if level is in config
	// if not, return
	if WARN&l.levels == 0 {
		return
	}
	e := Entry{Level: WARN, Message: msg}

	e.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e)
		l.runHook(e)
	} else {
		l.openEntry(e.enc)
	}

	l.closeEntry(e)
	l.finalizeIfContext(e)

	e.enc.Release()
}

// WarnWith returns a ChainEntry with WARN level
func (l *Logger) WarnWith(msg string) ChainEntry {
	// first find writer for level
	// if none, stop
	e := ChainEntry{
		Entry: Entry{
			l:       l,
			Level:   WARN,
			Message: msg,
		},
	}
	e.disabled = WARN&e.l.levels == 0
	if e.disabled {
		return e
	}

	e.Entry.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e.Entry)
		l.runHook(e.Entry)
		return e
	}

	l.openEntry(e.Entry.enc)
	return e
}

// WarnWithFields logs an entry with WARN level and custom fields.
func (l *Logger) WarnWithFields(msg string, fields func(Entry)) {
	if WARN&l.levels == 0 {
		return
	}
	e := Entry{
		Level:   WARN,
		Message: msg,
	}

	e.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e)
		l.runHook(e)
	} else {
		l.openEntry(e.enc)
	}

	fields(e)
	l.closeEntry(e)
	l.finalizeIfContext(e)

	e.enc.Release()
}

// Error logs an entry with ERROR level
func (l *Logger) Error(msg string) {
	if ERROR&l.levels == 0 {
		return
	}
	e := Entry{
		Level:   ERROR,
		Message: msg,
	}

	e.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e)
		l.runHook(e)
	} else {
		l.openEntry(e.enc)
	}

	l.closeEntry(e)
	l.finalizeIfContext(e)

	e.enc.Release()
}

// ErrorWith returns a ChainEntry with ERROR level.
func (l *Logger) ErrorWith(msg string) ChainEntry {
	// first find writer for level
	// if none, stop
	e := ChainEntry{
		Entry: Entry{
			l:       l,
			Level:   ERROR,
			Message: msg,
		},
	}
	e.disabled = ERROR&e.l.levels == 0
	if e.disabled {
		return e
	}

	e.Entry.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e.Entry)
		l.runHook(e.Entry)
		return e
	}

	l.openEntry(e.Entry.enc)
	return e
}

// ErrorWithFields logs an entry with ERROR level and custom fields.
func (l *Logger) ErrorWithFields(msg string, fields func(Entry)) {
	if ERROR&l.levels == 0 {
		return
	}
	e := Entry{
		Level:   ERROR,
		Message: msg,
	}

	e.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e)
		l.runHook(e)
	} else {
		l.openEntry(e.enc)
	}

	fields(e)
	l.closeEntry(e)
	l.finalizeIfContext(e)
	e.enc.Release()
}

// Fatal logs an entry with FATAL level.
func (l *Logger) Fatal(msg string) {
	if FATAL&l.levels == 0 {
		return
	}
	e := Entry{
		Level:   FATAL,
		Message: msg,
	}

	e.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e)
		l.runHook(e)
	} else {
		l.openEntry(e.enc)
	}

	l.closeEntry(e)
	l.finalizeIfContext(e)

	e.enc.Release()

	l.exit(1)
}

// FatalWith returns a ChainEntry with FATAL level.
func (l *Logger) FatalWith(msg string) ChainEntry {
	// first find writer for level
	// if none, stop
	e := ChainEntry{
		Entry: Entry{
			l:       l,
			Level:   FATAL,
			Message: msg,
		},
	}
	e.disabled = FATAL&e.l.levels == 0
	if e.disabled {
		return e
	}

	e.Entry.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e.Entry)
		l.runHook(e.Entry)
		return e
	}

	l.openEntry(e.Entry.enc)
	e.exit = true
	return e
}

// FatalWithFields logs an entry with FATAL level and custom fields.
func (l *Logger) FatalWithFields(msg string, fields func(Entry)) {
	if FATAL&l.levels == 0 {
		return
	}

	e := Entry{
		Level:   FATAL,
		Message: msg,
	}

	e.enc = gojay.BorrowEncoder(l.w)

	// if we do not require a context then we
	// format with formatter and return.
	if l.contextName == "" {
		l.beginEntry(e.Level, msg, e)
		l.runHook(e)
	} else {
		l.openEntry(e.enc)
	}

	fields(e)
	l.closeEntry(e)
	l.finalizeIfContext(e)

	e.enc.Release()
	l.exit(1)
}

func (l *Logger) openEntry(enc *Encoder) {
	enc.AppendBytes(logOpen)
}

func (l *Logger) beginEntry(level uint8, msg string, e Entry) {
	e.enc.AppendBytes(levelsJSON[level])
	e.enc.AppendString(msg)

	if l.ctx != nil && l.contextName == "" {
		for _, c := range l.ctx {
			c(e)
		}
	}
}

func (l Logger) runHook(e Entry) {
	if l.hook == nil {
		return
	}
	l.hook(e)
}

func (l *Logger) finalizeIfContext(entry Entry) {
	if l.contextName == "" {
		return
	}

	embeddedEnc := entry.enc

	// create a new encoder for the final output.
	entryEnc := gojay.BorrowEncoder(l.w)
	defer entryEnc.Release()

	entry.enc = entryEnc

	// create dummy entry for applying hooks.
	l.beginEntry(entry.Level, entry.Message, entry)
	l.runHook(entry)

	// Add entry's encoded data into new encoder.
	var embeddedJSON = gojay.EmbeddedJSON(embeddedEnc.Buf())
	entryEnc.AddEmbeddedJSONKey(l.contextName, &embeddedJSON)

	// close new encoder context for proper json.
	entryEnc.AppendBytes(logClose)

	// we need to manually write output as logger
	// has context.
	entryEnc.Write()
}

func (l *Logger) closeEntry(e Entry) {
	if l.contextName == "" {
		e.enc.AppendBytes(logClose)
	} else {
		if l.ctx != nil {
			for _, c := range l.ctx {
				c(e)
			}
		}
		e.enc.AppendBytes(logCloseOnly)
	}

	if l.contextName == "" {
		e.enc.Write()
	}
}

func (l *Logger) exit(code int) {
	if l.ExitFn == nil {
		// fallback to os.Exit to prevent panic incase set as nil.
		os.Exit(code)
	}
	l.ExitFn(code)
}

// Caller returns the caller in the stack trace, skipped n times.
func (l *Logger) Caller(n int) string {
	_, f, fl, _ := runtime.Caller(n)
	flStr := strconv.Itoa(fl)
	return f + ":" + flStr
}
