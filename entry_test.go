package onelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntry(t *testing.T) {
	t.Run("basic-info-entry", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.InfoWith("hello").Int("test", 1).Write()
		json := `{"level":"info","message":"hello","test":1}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-info-entry-disabled", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG)
		logger.InfoWith("hello").Int("test", 1).Write()
		json := ``
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-debug-entry", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.DebugWith("hello").Int("test", 1).Write()
		json := `{"level":"debug","message":"hello","test":1}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-debug-entry-disabled", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|FATAL)
		logger.DebugWith("hello").Int("test", 1).Write()
		json := ``
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-warn-entry", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.WarnWith("hello").Int("test", 1).Write()
		json := `{"level":"warn","message":"hello","test":1}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-warn-entry-disabled", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|ERROR|FATAL)
		logger.WarnWith("hello").Int("test", 1).Write()
		json := ``
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-error-entry", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.ErrorWith("hello").Int("test", 1).Write()
		json := `{"level":"error","message":"hello","test":1}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-error-entry-disabled", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|FATAL)
		logger.ErrorWith("hello").Int("test", 1).Write()
		json := ``
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-fatal-entry", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.FatalWith("hello").Int("test", 1).Write()
		json := `{"level":"fatal","message":"hello","test":1}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-fatal-entry-disabled", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR)
		logger.FatalWith("hello").Int("test", 1).Write()
		json := ``
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
}
