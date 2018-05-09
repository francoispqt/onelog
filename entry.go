package onelog

import "github.com/francoispqt/gojay"

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
	e.enc.AddStringKey(k, v)
	return e
}

// Int adds an int to the log entry.
func (e Entry) Int(k string, v int) Entry {
	e.enc.AddIntKey(k, v)
	return e
}

// Int64 adds an int64 to the log entry.
func (e Entry) Int64(k string, v int64) Entry {
	e.enc.AddInt64Key(k, v)
	return e
}

// Float adds a float64 to the log entry.
func (e Entry) Float(k string, v float64) Entry {
	e.enc.AddFloatKey(k, v)
	return e
}

// Bool adds a bool to the log entry.
func (e Entry) Bool(k string, v bool) Entry {
	e.enc.AddBoolKey(k, v)
	return e
}

// Error adds an error to the log entry.
func (e Entry) Err(k string, v error) Entry {
	e.enc.AddStringKey(k, v.Error())
	return e
}

// ObjectFunc adds an object to the log entry by calling a function.
func (e Entry) ObjectFunc(k string, v func(Entry)) Entry {
	e.enc.AddObjectKey(k, Object(func(enc *Encoder) {
		v(e)
	}))
	return e
}

// Object adds an object to the log entry by passing an implementation of gojay.MarshalerObject.
func (e Entry) Object(k string, obj gojay.MarshalerObject) Entry {
	e.enc.AddObjectKey(k, obj)
	return e
}

// Array adds an object to the log entry by passing an implementation of gojay.MarshalerObject.
func (e Entry) Array(k string, obj gojay.MarshalerArray) Entry {
	e.enc.AddArrayKey(k, obj)
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
	e.enc.AddStringKey(k, v)
	return e
}

// Int adds an int to the log entry.
func (e entry) Int(k string, v int) entry {
	if e.disabled {
		return e
	}
	e.enc.AddIntKey(k, v)
	return e
}

// Int64 adds an int64 to the log entry.
func (e entry) Int64(k string, v int64) entry {
	if e.disabled {
		return e
	}
	e.enc.AddInt64Key(k, v)
	return e
}

// Float adds a float64 to the log entry.
func (e entry) Float(k string, v float64) entry {
	if e.disabled {
		return e
	}
	e.enc.AddFloatKey(k, v)
	return e
}

// Bool adds a bool to the log entry.
func (e entry) Bool(k string, v bool) entry {
	if e.disabled {
		return e
	}
	e.enc.AddBoolKey(k, v)
	return e
}

// Error adds an error to the log entry.
func (e entry) Err(k string, v error) entry {
	if e.disabled {
		return e
	}
	e.enc.AddStringKey(k, v.Error())
	return e
}

// ObjectFunc adds an object to the log entry by calling a function.
func (e entry) ObjectFunc(k string, v func(Entry)) entry {
	if e.disabled {
		return e
	}
	e.enc.AddObjectKey(k, Object(func(enc *Encoder) {
		v(e.Entry)
	}))
	return e
}

// Object adds an object to the log entry by passing an implementation of gojay.MarshalerObject.
func (e entry) Object(k string, obj gojay.MarshalerObject) entry {
	if e.disabled {
		return e
	}
	e.enc.AddObjectKey(k, obj)
	return e
}

// Array adds an object to the log entry by passing an implementation of gojay.MarshalerObject.
func (e entry) Array(k string, obj gojay.MarshalerArray) entry {
	if e.disabled {
		return e
	}
	e.enc.AddArrayKey(k, obj)
	return e
}
