package onelog

import "github.com/francoispqt/gojay"

// Entry is the structure wrapping a pointer to the current encoder.
// It provides easy API to work with GoJay's encoder.
type Entry struct {
	enc *Encoder
}

// String adds a string to the log entry.
func (e Entry) String(k, v string) {
	e.enc.AddStringKey(k, v)
}

// Int adds an int to the log entry.
func (e Entry) Int(k string, v int) {
	e.enc.AddIntKey(k, v)
}

// Int adds an int to the log entry.
func (e Entry) Int64(k string, v int64) {
	e.enc.AddInt64Key(k, v)
}

// Bool adds a bool to the log entry.
func (e Entry) Bool(k string, v bool) {
	e.enc.AddBoolKey(k, v)
}

// Error adds an error to the log entry.
func (e Entry) Error(k string, v error) {
	e.enc.AddStringKey(k, v.Error())
}

// ObjectFunc adds an object to the log entry by calling a function.
func (e Entry) ObjectFunc(k string, v func()) {
	e.enc.AddObjectKey(k, Object(func(enc *Encoder) {
		v()
	}))
}

// Object adds an object to the log entry by passing an implementation of gojay.MarshalerObject.
func (e Entry) Object(k string, obj gojay.MarshalerObject) {
	e.enc.AddObjectKey(k, obj)
}

// Array adds an object to the log entry by passing an implementation of gojay.MarshalerObject.
func (e Entry) Array(k string, obj gojay.MarshalerArray) {
	e.enc.AddArrayKey(k, obj)
}
