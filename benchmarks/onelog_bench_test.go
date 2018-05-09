package benchmarks

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/francoispqt/onelog"
)

func TestPrint(t *testing.T) {
	logger := onelog.New(os.Stdout, onelog.ALL).
		Hook(func(e onelog.Entry) {
			e.Int64("time", time.Now().Unix())
		})
	logger.InfoWith("message").
		String("test", "test").
		String("test", "test").
		String("test", "test").
		String("test", "test").
		String("test", "test").
		String("test", "test").
		String("test", "test").
		Write()
}

func BenchmarkOnelog(b *testing.B) {
	b.Run("with-fields", func(b *testing.B) {
		logger := onelog.New(ioutil.Discard, onelog.ALL).
			Hook(func(e onelog.Entry) {
				e.Int64("time", time.Now().Unix())
			})
		s := struct {
			i int
		}{i: 0}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoWithFields("message", func(e onelog.Entry) {
					e.Int("test", s.i)
					e.String("test", "test")
					e.String("test", "test")
					e.String("test", "test")
					e.String("test", "test")
					e.String("test", "test")
					e.String("test", "test")
					e.String("test", "test")
				})
			}
		})
	})
	b.Run("message-only", func(b *testing.B) {
		logger := onelog.New(ioutil.Discard, onelog.ALL).
			Hook(func(e onelog.Entry) {
				e.Int64("time", time.Now().Unix())
			})
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("message")
			}
		})
	})
	b.Run("entry-message-only", func(b *testing.B) {
		logger := onelog.New(ioutil.Discard, onelog.ALL).
			Hook(func(e onelog.Entry) {
				e.Int64("time", time.Now().Unix())
			})
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoWith("message").Write()
			}
		})
	})
	b.Run("entry-fields", func(b *testing.B) {
		logger := onelog.New(ioutil.Discard, onelog.ALL).
			Hook(func(e onelog.Entry) {
				e.Int64("time", time.Now().Unix())
			})
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoWith("message").
					String("test", "test").
					String("test", "test").
					String("test", "test").
					String("test", "test").
					String("test", "test").
					String("test", "test").
					String("test", "test").
					Write()
			}
		})
	})
}
