package benchmarks

import (
	"io/ioutil"
	"testing"

	"github.com/francoispqt/onelog"
)

func BenchmarkOnelog(b *testing.B) {
	b.Run("message-only", func(b *testing.B) {
		Logger := onelog.NewLogger(ioutil.Discard, onelog.INFO|onelog.WARN|onelog.FATAL|onelog.ERROR)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				Logger.Info("message")
			}
		})
	})
	b.Run("with-fields", func(b *testing.B) {
		Logger := onelog.NewLogger(ioutil.Discard, onelog.INFO|onelog.WARN|onelog.FATAL|onelog.ERROR) // for non blocking NewStreamLogger
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				Logger.InfoWithFields("message", func(enc *onelog.Encoder) {
					enc.AddStringKey("test", "test")
					enc.AddStringKey("test", "test")
					enc.AddStringKey("test", "test")
					enc.AddStringKey("test", "test")
					enc.AddStringKey("test", "test")
					enc.AddStringKey("test", "test")
					enc.AddStringKey("test", "test")
				})
			}
		})
	})
}
