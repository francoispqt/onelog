package onelog

import (
	"github.com/francoispqt/gojay"
)

type entry struct {
	Entry
	disabled bool
}

// Entry is the structure wrapping a pointer to the current encoder.
// It provides easy API to work with GoJay's encoder.
type Entry struct {
	enc *Encoder
	l   *Logger
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
	e.enc.StringKey(k, v.Error())
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

// entry it is used when chaining

// Info logs an entry with INFO level.
func (e entry) Write() {
	if e.disabled {
		return
	}
	// first find writer for level
	// if none, stop
	e.Entry.l.closeEntry(e.enc)
	e.Entry.enc.Release()
}

// String adds a string to the log entry.
func (e entry) String(k, v string) entry {
	if e.disabled {
		return e
	}
	e.enc.StringKey(k, v)
	return e
}

// Int adds an int to the log entry.
func (e entry) Int(k string, v int) entry {
	if e.disabled {
		return e
	}
	e.enc.IntKey(k, v)
	return e
}

// Int64 adds an int64 to the log entry.
func (e entry) Int64(k string, v int64) entry {
	if e.disabled {
		return e
	}
	e.enc.Int64Key(k, v)
	return e
}

// Float adds a float64 to the log entry.
func (e entry) Float(k string, v float64) entry {
	if e.disabled {
		return e
	}
	e.enc.FloatKey(k, v)
	return e
}

// Bool adds a bool to the log entry.
func (e entry) Bool(k string, v bool) entry {
	if e.disabled {
		return e
	}
	e.enc.BoolKey(k, v)
	return e
}

// Error adds an error to the log entry.
func (e entry) Err(k string, v error) entry {
	if e.disabled {
		return e
	}
	e.enc.StringKey(k, v.Error())
	return e
}

// ObjectFunc adds an object to the log entry by calling a function.
func (e entry) ObjectFunc(k string, v func(Entry)) entry {
	if e.disabled {
		return e
	}
	e.enc.ObjectKey(k, Object(func(enc *Encoder) {
		v(e.Entry)
	}))
	return e
}

// Object adds an object to the log entry by passing an implementation of gojay.MarshalerJSONObject.
func (e entry) Object(k string, obj gojay.MarshalerJSONObject) entry {
	if e.disabled {
		return e
	}
	e.enc.ObjectKey(k, obj)
	return e
}

// Array adds an object to the log entry by passing an implementation of gojay.MarshalerJSONObject.
func (e entry) Array(k string, obj gojay.MarshalerJSONArray) entry {
	if e.disabled {
		return e
	}
	e.enc.ArrayKey(k, obj)
	return e
}

// Any adds anything stuff to the log entry based on it's type
func (e entry) Any(k string, obj interface{}) entry {
	if e.disabled {
		return e
	}
	e.enc.AddInterfaceKey(k, obj)
	return e
}
