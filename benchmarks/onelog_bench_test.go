package benchmarks

import (
	"io/ioutil"
	"testing"

	"github.com/francoispqt/onelog"
)

func BenchmarkOnelog(b *testing.B) {
	b.Run("with-fields", func(b *testing.B) {
		logger := onelog.NewLogger(ioutil.Discard, onelog.ALL) // for non blocking NewStreamLogger
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.InfoWithFields("message", func(e onelog.Entry) {
					e.String("test", "test")
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
		logger := onelog.NewLogger(ioutil.Discard, onelog.ALL)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("message")
			}
		})
	})
}
