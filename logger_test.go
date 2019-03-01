package onelog

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/francoispqt/gojay"
	"github.com/stretchr/testify/assert"
)

type TestWriter struct {
	b      []byte
	called bool
}

func (t *TestWriter) Write(b []byte) (int, error) {
	t.called = true
	if len(t.b) < len(b) {
		t.b = make([]byte, len(b))
	}
	copy(t.b, b)
	return len(t.b), nil
}

func newWriter() *TestWriter {
	return &TestWriter{make([]byte, 0, 512), false}
}

type TestObj struct {
	foo string
}

func (t *TestObj) MarshalJSONObject(enc *gojay.Encoder) {
	enc.AddStringKey("foo", t.foo)
}

func (t *TestObj) IsNil() bool {
	return t == nil
}

type TestObjArr []*TestObj

func (t TestObjArr) MarshalJSONArray(enc *gojay.Encoder) {
	for _, o := range t {
		enc.AddObject(o)
	}
}

func (t TestObjArr) IsNil() bool {
	return len(t) == 0
}

func TestOnelogFeature(t *testing.T) {
	t.Run("custom-msg-key", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		MsgKey("test")
		logger.Info("message")
		assert.Equal(t, `{"level":"info","test":"message"}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
		MsgKey("message")
	})
	t.Run("custom-lvl-key", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		LevelKey("test")
		logger.Info("message")
		assert.Equal(t, `{"test":"info","message":"message"}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
		LevelKey("level")
	})
	t.Run("custom-lvl-text", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		LevelText(DEBUG, "DEBUG")
		logger.Debug("message")
		assert.Equal(t, `{"level":"DEBUG","message":"message"}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
		LevelText(DEBUG, "debug")
	})
	t.Run("caller", func(t *testing.T) {
		logger := New(nil, DEBUG|INFO|WARN|ERROR|FATAL)
		str := logger.Caller(1)
		strs := strings.Split(str, "/")
		assert.Equal(t, "logger_test.go:83", strs[len(strs)-1], "file should be logger_test.go:81")
	})
}
func TestOnelogWithoutFields(t *testing.T) {
	t.Run("basic-message-info", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Info("message")
		assert.Equal(t, `{"level":"info","message":"message"}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-debug", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Debug("message")
		assert.Equal(t, `{"level":"debug","message":"message"}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-warn", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Warn("message")
		assert.Equal(t, `{"level":"warn","message":"message"}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-error", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Error("message")
		assert.Equal(t, `{"level":"error","message":"message"}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-fatal", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			assert.Equal(t, `{"level":"fatal","message":"message"}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.Fatal("message")
	})
	t.Run("basic-message-disabled-level-info", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|WARN|ERROR|FATAL)
		logger.Info("message")
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-disabled-level-debug", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|FATAL)
		logger.Debug("message")
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-disabled-level-warn", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|DEBUG|ERROR|FATAL)
		logger.Warn("message")
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-disabled-level-error", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|DEBUG|FATAL)
		logger.Error("message")
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-disabled-level-fatal", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|DEBUG)
		logger.Fatal("message")
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
	})
}

func TestOnelogContextWithoutFields(t *testing.T) {
	t.Run("basic-message-info", func(t *testing.T) {
		w := newWriter()
		logger := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger.Info("message")
		assert.Equal(t, `{"level":"info","message":"message","params":{}}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-debug", func(t *testing.T) {
		w := newWriter()
		logger := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger.Debug("message")
		assert.Equal(t, `{"level":"debug","message":"message","params":{}}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-warn", func(t *testing.T) {
		w := newWriter()
		logger := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger.Warn("message")
		assert.Equal(t, `{"level":"warn","message":"message","params":{}}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-error", func(t *testing.T) {
		w := newWriter()
		logger := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger.Error("message")
		assert.Equal(t, `{"level":"error","message":"message","params":{}}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-fatal", func(t *testing.T) {
		w := newWriter()
		logger := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			assert.Equal(t, `{"level":"fatal","message":"message","params":{}}`+"\n", string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.Fatal("message")
	})
	t.Run("basic-message-disabled-level-info", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|WARN|ERROR|FATAL)
		logger.Info("message")
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-disabled-level-debug", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|FATAL)
		logger.Debug("message")
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-disabled-level-warn", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|DEBUG|ERROR|FATAL)
		logger.Warn("message")
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-disabled-level-error", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|DEBUG|FATAL)
		logger.Error("message")
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
	})
	t.Run("basic-message-disabled-level-fatal", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|DEBUG)
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
			r := recover()
			if r != nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.Fatal("message")
	})
}

func TestOnelogWithFields(t *testing.T) {
	t.Run("fields-info", func(t *testing.T) {
		testObj := &TestObj{foo: "bar"}
		testArr := TestObjArr{testObj}
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.InfoWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int("count", 100)
			e.Int64("int64", 100)
			e.Float("float64", 0.15)
			e.Bool("done", true)
			e.Err("error", errors.New("some error"))
			e.ObjectFunc("user", func(e Entry) {
				e.String("name", "somename")
			})
			e.Object("testObj", testObj)
			e.Array("testArr", testArr)
		})
		json := `{"level":"info","message":"message","userID":"123456",` +
			`"action":"login","result":"success","count":100,"int64":100,"float64":0.15,"done":true,` +
			`"error":"some error","user":{"name":"somename"},"testObj":{"foo":"bar"},` +
			`"testArr":[{"foo":"bar"}]}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-debug", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.DebugWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"debug","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-warn", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.WarnWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"warn","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-error", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.ErrorWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"error","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-fatal", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"message","userID":"123456","action":"login","result":"success","int64":120}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.FatalWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int64("int64", 120)
		})
	})
	t.Run("fields-disabled-level-info", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|WARN|ERROR|FATAL)
		logger.InfoWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-debug", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|FATAL)
		logger.DebugWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-warn", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|DEBUG|ERROR|FATAL)
		logger.WarnWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-error", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|DEBUG|FATAL)
		logger.ErrorWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-fatal", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|DEBUG)
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
			assert.False(t, w.called, "writer should not be called")
			r := recover()
			if r != nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.FatalWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
	})
}

func TestOnelogHook(t *testing.T) {
	t.Run("hook-basic-info", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Hook(func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		logger.Info("message")
		json := `{"level":"info","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("hook-basic-debug", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Hook(func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		logger.Debug("message")
		json := `{"level":"debug","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("hook-basic-warn", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Hook(func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		logger.Warn("message")
		json := `{"level":"warn","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("hook-basic-error", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Hook(func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		logger.Error("message")
		json := `{"level":"error","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("hook-basic-fatal", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Hook(func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.Fatal("message")
	})
	t.Run("hook-fields-info", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Hook(func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		logger.InfoWithFields("message", func(e Entry) {
			e.String("field", "field")
		})
		json := `{"level":"info","message":"message","userID":"123456","action":"login","result":"success","field":"field"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("hook-fields-debug", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Hook(func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		logger.DebugWithFields("message", func(e Entry) {
			e.String("field", "field")
		})
		json := `{"level":"debug","message":"message","userID":"123456","action":"login","result":"success","field":"field"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("hook-fields-warn", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Hook(func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		logger.WarnWithFields("message", func(e Entry) {
			e.String("field", "field")
		})
		json := `{"level":"warn","message":"message","userID":"123456","action":"login","result":"success","field":"field"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("hook-fields-error", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.Hook(func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		logger.ErrorWithFields("message", func(e Entry) {
			e.String("field", "field")
		})
		json := `{"level":"error","message":"message","userID":"123456","action":"login","result":"success","field":"field"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("hook-fields-fatal", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|INFO|WARN|ERROR|FATAL)
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"message","userID":"123456","action":"login","result":"success","field":"field"}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.Hook(func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		logger.FatalWithFields("message", func(e Entry) {
			e.String("field", "field")
		})
	})
}

func TestOnelogContext(t *testing.T) {
	t.Run("context-info-basic", func(t *testing.T) {
		w := newWriter()
		logger := New(w, ALL).With(func(e Entry) {
			e.String("test", "test")
		})
		logger.Info("test")
		json := `{"level":"info","message":"test","test":"test"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")

	})
	t.Run("context-info-fields", func(t *testing.T) {
		w := newWriter()
		logger := New(w, ALL).With(func(e Entry) {
			e.String("test", "test")
		})
		logger.InfoWithFields("test", func(e Entry) {
			e.String("field", "field")
		})
		json := `{"level":"info","message":"test","test":"test","field":"field"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})

	t.Run("context-debug-basic", func(t *testing.T) {
		w := newWriter()
		logger := New(w, ALL).With(func(e Entry) {
			e.String("test", "test")
		})
		logger.Debug("test")
		json := `{"level":"debug","message":"test","test":"test"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")

	})
	t.Run("context-debug-fields", func(t *testing.T) {
		w := newWriter()
		logger := New(w, ALL).With(func(e Entry) {
			e.String("test", "test")
		})
		logger.DebugWithFields("test", func(e Entry) {
			e.String("field", "field")
		})
		json := `{"level":"debug","message":"test","test":"test","field":"field"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})

	t.Run("context-warn-basic", func(t *testing.T) {
		w := newWriter()
		logger := New(w, ALL).With(func(e Entry) {
			e.String("test", "test")
		})
		logger.Warn("test")
		json := `{"level":"warn","message":"test","test":"test"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")

	})
	t.Run("context-warn-fields", func(t *testing.T) {
		w := newWriter()
		logger := New(w, ALL).With(func(e Entry) {
			e.String("test", "test")
		})
		logger.WarnWithFields("test", func(e Entry) {
			e.String("field", "field")
		})
		json := `{"level":"warn","message":"test","test":"test","field":"field"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})

	t.Run("context-error-basic", func(t *testing.T) {
		w := newWriter()
		logger := New(w, ALL).With(func(e Entry) {
			e.String("test", "test")
		})
		logger.Error("test")
		json := `{"level":"error","message":"test","test":"test"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")

	})
	t.Run("context-error-fields", func(t *testing.T) {
		w := newWriter()
		logger := New(w, ALL).With(func(e Entry) {
			e.String("test", "test")
		})
		logger.ErrorWithFields("test", func(e Entry) {
			e.String("field", "field")
		})
		json := `{"level":"error","message":"test","test":"test","field":"field"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})

	t.Run("context-fatal-basic", func(t *testing.T) {
		w := newWriter()
		logger := New(w, ALL).With(func(e Entry) {
			e.String("test", "test")
		})
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"test","test":"test"}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.Fatal("test")
	})
	t.Run("context-fatal-fields", func(t *testing.T) {
		w := newWriter()
		logger := New(w, ALL).With(func(e Entry) {
			e.String("test", "test")
		})
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"test","test":"test","field":"field"}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.FatalWithFields("test", func(e Entry) {
			e.String("field", "field")
		})
	})
}

func TestNilOutput(t *testing.T) {
	t.Run("nil-logger", func(t *testing.T) {
		logger := New(nil, 0)
		assert.NotNil(t, logger.w, "Logger output should not be nil")
		logger.Error("Test output")
		t.Log("Successfully write to ioutil.Discard output")
	})
	t.Run("nil-context-logger", func(t *testing.T) {
		logger := NewContext(nil, 0, "params")
		assert.NotNil(t, logger.w, "Logger output should not be nil")
		logger.Error("Test output")
		t.Log("Successfully write to ioutil.Discard output")
	})
}

func TestOnelogHooksWithAndContext(t *testing.T) {
	t.Run("fields-info", func(t *testing.T) {
		testObj := &TestObj{foo: "bar"}
		testArr := TestObjArr{testObj}
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.Hook(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.InfoWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int("count", 100)
			e.Int64("int64", 100)
			e.Float("float64", 0.15)
			e.Bool("done", true)
			e.Err("error", errors.New("some error"))
			e.ObjectFunc("user", func(e Entry) {
				e.String("name", "somename")
			})
			e.Object("testObj", testObj)
			e.Array("testArr", testArr)
		})
		json := `{"level":"info","message":"message","thunder_frequency":1000,"params":{"userID":"123456",` +
			`"action":"login","result":"success","count":100,"int64":100,"float64":0.15,"done":true,` +
			`"error":"some error","user":{"name":"somename"},"testObj":{"foo":"bar"},` +
			`"testArr":[{"foo":"bar"}]}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-debug", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.Hook(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.DebugWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"debug","message":"message","thunder_frequency":1000,"params":{"userID":"123456","action":"login","result":"success"}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-warn", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.Hook(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.WarnWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"warn","message":"message","thunder_frequency":1000,"params":{"userID":"123456","action":"login","result":"success"}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-error", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.Hook(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.ErrorWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"error","message":"message","thunder_frequency":1000,"params":{"userID":"123456","action":"login","result":"success"}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-fatal", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.Hook(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"message","thunder_frequency":1000,"params":{"userID":"123456","action":"login","result":"success","int64":120}}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.FatalWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int64("int64", 120)
		})
	})
}

func TestOnelogWithAndContext(t *testing.T) {
	t.Run("fields-info", func(t *testing.T) {
		testObj := &TestObj{foo: "bar"}
		testArr := TestObjArr{testObj}
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.InfoWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int("count", 100)
			e.Int64("int64", 100)
			e.Float("float64", 0.15)
			e.Bool("done", true)
			e.Err("error", errors.New("some error"))
			e.ObjectFunc("user", func(e Entry) {
				e.String("name", "somename")
			})
			e.Object("testObj", testObj)
			e.Array("testArr", testArr)
		})
		json := `{"level":"info","message":"message","params":{"userID":"123456",` +
			`"action":"login","result":"success","count":100,"int64":100,"float64":0.15,"done":true,` +
			`"error":"some error","user":{"name":"somename"},"testObj":{"foo":"bar"},` +
			`"testArr":[{"foo":"bar"}],"thunder_frequency":1000}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-debug", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.DebugWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"debug","message":"message","params":{"userID":"123456","action":"login","result":"success","thunder_frequency":1000}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-warn", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.WarnWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"warn","message":"message","params":{"userID":"123456","action":"login","result":"success","thunder_frequency":1000}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-error", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.ErrorWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"error","message":"message","params":{"userID":"123456","action":"login","result":"success","thunder_frequency":1000}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-fatal", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"message","params":{"userID":"123456","action":"login","result":"success","int64":120,"thunder_frequency":1000}}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.FatalWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int64("int64", 120)
		})
	})
}

func TestOnelogWithFieldsAndContext(t *testing.T) {
	t.Run("fields-info", func(t *testing.T) {
		testObj := &TestObj{foo: "bar"}
		testArr := TestObjArr{testObj}
		w := newWriter()
		logger := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger.InfoWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int("count", 100)
			e.Int64("int64", 100)
			e.Float("float64", 0.15)
			e.Bool("done", true)
			e.Err("error", errors.New("some error"))
			e.ObjectFunc("user", func(e Entry) {
				e.String("name", "somename")
			})
			e.Object("testObj", testObj)
			e.Array("testArr", testArr)
		})
		json := `{"level":"info","message":"message","params":{"userID":"123456",` +
			`"action":"login","result":"success","count":100,"int64":100,"float64":0.15,"done":true,` +
			`"error":"some error","user":{"name":"somename"},"testObj":{"foo":"bar"},` +
			`"testArr":[{"foo":"bar"}]}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-debug", func(t *testing.T) {
		w := newWriter()
		logger := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger.DebugWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"debug","message":"message","params":{"userID":"123456","action":"login","result":"success"}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-warn", func(t *testing.T) {
		w := newWriter()
		logger := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger.WarnWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"warn","message":"message","params":{"userID":"123456","action":"login","result":"success"}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-error", func(t *testing.T) {
		w := newWriter()
		logger := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger.ErrorWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"error","message":"message","params":{"userID":"123456","action":"login","result":"success"}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-fatal", func(t *testing.T) {
		w := newWriter()
		logger := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"message","params":{"userID":"123456","action":"login","result":"success","int64":120}}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.FatalWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int64("int64", 120)
		})
	})
	t.Run("fields-disabled-level-info", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|WARN|ERROR|FATAL)
		logger.InfoWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-debug", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|FATAL)
		logger.DebugWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-warn", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|DEBUG|ERROR|FATAL)
		logger.WarnWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-error", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|DEBUG|FATAL)
		logger.ErrorWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-fatal", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|DEBUG)
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
			assert.False(t, w.called, "writer should not be called")
			r := recover()
			if r != nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.FatalWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
	})
}

func TestOnelogNoContextFromContextLogger(t *testing.T) {
	t.Run("fields-info", func(t *testing.T) {
		testObj := &TestObj{foo: "bar"}
		testArr := TestObjArr{testObj}
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.WithContext("")
		logger.InfoWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int("count", 100)
			e.Int64("int64", 100)
			e.Float("float64", 0.15)
			e.Bool("done", true)
			e.Err("error", errors.New("some error"))
			e.ObjectFunc("user", func(e Entry) {
				e.String("name", "somename")
			})
			e.Object("testObj", testObj)
			e.Array("testArr", testArr)
		})
		json := `{"level":"info","message":"message","userID":"123456",` +
			`"action":"login","result":"success","count":100,"int64":100,"float64":0.15,"done":true,` +
			`"error":"some error","user":{"name":"somename"},"testObj":{"foo":"bar"},` +
			`"testArr":[{"foo":"bar"}]}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-debug", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.WithContext("")
		logger.DebugWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"debug","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-warn", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.WithContext("")
		logger.WarnWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"warn","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-error", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.WithContext("")
		logger.ErrorWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		json := `{"level":"error","message":"message","userID":"123456","action":"login","result":"success"}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-fatal", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.WithContext("")
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"message","userID":"123456","action":"login","result":"success","int64":120}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.FatalWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int64("int64", 120)
		})
	})
	t.Run("fields-disabled-level-info", func(t *testing.T) {
		w := newWriter()
		logger := New(w, DEBUG|WARN|ERROR|FATAL)
		logger.InfoWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-debug", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|FATAL)
		logger.DebugWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-warn", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|DEBUG|ERROR|FATAL)
		logger.WarnWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-error", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|DEBUG|FATAL)
		logger.ErrorWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
		assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
		assert.False(t, w.called, "writer should not be called")
	})
	t.Run("basic-message-disabled-level-fatal", func(t *testing.T) {
		w := newWriter()
		logger := New(w, INFO|WARN|ERROR|DEBUG)
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			assert.Equal(t, string(w.b), ``, "bytes written to the writer dont equal expected result")
			assert.False(t, w.called, "writer should not be called")
			r := recover()
			if r != nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.FatalWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
		})
	})
}

func TestOnelogFieldsChainAndContext(t *testing.T) {
	t.Run("fields-info", func(t *testing.T) {
		testObj := &TestObj{foo: "bar"}
		testArr := TestObjArr{testObj}
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.InfoWith("message").
			String("userID", "123456").
			String("action", "login").
			String("result", "success").
			Int("count", 100).
			Int64("int64", 100).
			Float("float64", 0.15).
			Bool("done", true).
			Err("error", errors.New("some error")).
			ObjectFunc("user", func(e Entry) {
				e.String("name", "somename")
			}).
			Object("testObj", testObj).
			Array("testArr", testArr).Write()

		json := `{"level":"info","message":"message","params":{"userID":"123456",` +
			`"action":"login","result":"success","count":100,"int64":100,"float64":0.15,"done":true,` +
			`"error":"some error","user":{"name":"somename"},"testObj":{"foo":"bar"},` +
			`"testArr":[{"foo":"bar"}],"thunder_frequency":1000}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-debug", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.DebugWith("message").
			String("userID", "123456").
			String("action", "login").
			String("result", "success").Write()
		json := `{"level":"debug","message":"message","params":{"userID":"123456","action":"login","result":"success","thunder_frequency":1000}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-warn", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.WarnWith("message").
			String("userID", "123456").
			String("action", "login").
			String("result", "success").Write()
		json := `{"level":"warn","message":"message","params":{"userID":"123456","action":"login","result":"success","thunder_frequency":1000}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-error", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.ErrorWith("message").
			String("userID", "123456").
			String("action", "login").
			String("result", "success").Write()
		json := `{"level":"error","message":"message","params":{"userID":"123456","action":"login","result":"success","thunder_frequency":1000}}` + "\n"
		assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
	})
	t.Run("fields-fatal", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"message","params":{"userID":"123456","action":"login","result":"success","int64":120,"thunder_frequency":1000}}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.FatalWith("message").
			String("userID", "123456").
			String("action", "login").
			String("result", "success").
			Int64("int64", 120).Write()

	})

	t.Run("multi-augmented-logger", func(t *testing.T) {
		w := newWriter()
		parent := NewContext(w, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.With(func(e Entry) { e.Int("thunder_frequency", 1000) })
		logger = logger.With(func(e Entry) { e.String("foo", "bar") })
		logger.ExitFn = func(c int) {
			panic("os.Exit called")
		}
		defer func() {
			json := `{"level":"fatal","message":"message","params":{"thunder_frequency":1000,"foo":"bar"}}` + "\n"
			assert.Equal(t, json, string(w.b), "bytes written to the writer dont equal expected result")
			r := recover()
			if r == nil {
				t.Errorf("logger.Fatal() recover = %v", r)
			}
		}()
		logger.Fatal("message")
	})
}

func TestFatalActualOsExit(t *testing.T) {
	if os.Getenv("FatalActualOsExit") == "1" {
		parent := NewContext(os.Stdout, DEBUG|INFO|WARN|ERROR|FATAL, "params")
		logger := parent.WithContext("")
		logger.ExitFn = nil
		logger.FatalWithFields("message", func(e Entry) {
			e.String("userID", "123456")
			e.String("action", "login")
			e.String("result", "success")
			e.Int64("int64", 120)
		})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalActualOsExit")
	cmd.Env = append(os.Environ(), "FatalActualOsExit=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
