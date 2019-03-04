package onelog

import (
	"github.com/francoispqt/gojay"
)

// Entry is the structure wrapping a pointer to the current encoder.
// It provides easy API to work with GoJay's encoder.
type Entry struct {
	enc     *Encoder
	l       *Logger
	Level   uint8
	Message string
}

// String adds a string to the log entry.
func (e Entry) String(k, v string) Entry {
	e.enc.StringKey(k, v)
	return e
}

// Int adds an int to the log entry.
func (e Entry) Int(k string, v int) Entry {
	e.enc.IntKey(k, v)
	return e
}

// Int64 adds an int64 to the log entry.
func (e Entry) Int64(k string, v int64) Entry {
	e.enc.Int64Key(k, v)
	return e
}

// Float adds a float64 to the log entry.
func (e Entry) Float(k string, v float64) Entry {
	e.enc.FloatKey(k, v)
	return e
}

// Bool adds a bool to the log entry.
func (e Entry) Bool(k string, v bool) Entry {
	e.enc.BoolKey(k, v)
	return e
}

// Err adds an error to the log entry.
func (e Entry) Err(k string, v error) Entry {
	if v != nil {
		e.enc.StringKey(k, v.Error())
	}
	return e
}

// ObjectFunc adds an object to the log entry by calling a function.
func (e Entry) ObjectFunc(k string, v func(Entry)) Entry {
	e.enc.ObjectKey(k, Object(func(enc *Encoder) {
		v(e)
	}))
	return e
}

// Object adds an object to the log entry by passing an implementation of gojay.MarshalerJSONObject.
func (e Entry) Object(k string, obj gojay.MarshalerJSONObject) Entry {
	e.enc.ObjectKey(k, obj)
	return e
}

// Array adds an object to the log entry by passing an implementation of gojay.MarshalerJSONObject.
func (e Entry) Array(k string, obj gojay.MarshalerJSONArray) Entry {
	e.enc.ArrayKey(k, obj)
	return e
}

// ChainEntry is for chaining calls to the entry.
type ChainEntry struct {
	Entry
	disabled bool
	exit     bool
}

// Info logs an entry with INFO level.
func (e ChainEntry) Write() {
	if e.disabled {
		return
	}
	// first find writer for level
	// if none, stop
	e.Entry.l.closeEntry(e.Entry)
	e.Entry.l.finalizeIfContext(e.Entry)
	e.Entry.enc.Release()

	if e.exit {
		e.Entry.l.exit(1)
	}
}

// String adds a string to the log entry.
func (e ChainEntry) String(k, v string) ChainEntry {
	if e.disabled {
		return e
	}
	e.enc.StringKey(k, v)
	return e
}

// Int adds an int to the log entry.
func (e ChainEntry) Int(k string, v int) ChainEntry {
	if e.disabled {
		return e
	}
	e.enc.IntKey(k, v)
	return e
}

// Int64 adds an int64 to the log entry.
func (e ChainEntry) Int64(k string, v int64) ChainEntry {
	if e.disabled {
		return e
	}
	e.enc.Int64Key(k, v)
	return e
}

// Float adds a float64 to the log entry.
func (e ChainEntry) Float(k string, v float64) ChainEntry {
	if e.disabled {
		return e
	}
	e.enc.FloatKey(k, v)
	return e
}

// Bool adds a bool to the log entry.
func (e ChainEntry) Bool(k string, v bool) ChainEntry {
	if e.disabled {
		return e
	}
	e.enc.BoolKey(k, v)
	return e
}

// Err adds an error to the log entry.
func (e ChainEntry) Err(k string, v error) ChainEntry {
	if e.disabled {
		return e
	}
	if v != nil {
		e.enc.StringKey(k, v.Error())
	}
	return e
}

// ObjectFunc adds an object to the log entry by calling a function.
func (e ChainEntry) ObjectFunc(k string, v func(Entry)) ChainEntry {
	if e.disabled {
		return e
	}
	e.enc.ObjectKey(k, Object(func(enc *Encoder) {
		v(e.Entry)
	}))
	return e
}

// Object adds an object to the log entry by passing an implementation of gojay.MarshalerJSONObject.
func (e ChainEntry) Object(k string, obj gojay.MarshalerJSONObject) ChainEntry {
	if e.disabled {
		return e
	}
	e.enc.ObjectKey(k, obj)
	return e
}

// Array adds an object to the log entry by passing an implementation of gojay.MarshalerJSONObject.
func (e ChainEntry) Array(k string, obj gojay.MarshalerJSONArray) ChainEntry {
	if e.disabled {
		return e
	}
	e.enc.ArrayKey(k, obj)
	return e
}

// Any adds anything stuff to the log entry based on it's type
func (e ChainEntry) Any(k string, obj interface{}) ChainEntry {
	if e.disabled {
		return e
	}
	e.enc.AddInterfaceKey(k, obj)
	return e
}
