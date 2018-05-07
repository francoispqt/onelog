// Package onelog is a fast, low allocation and modular JSON logger.
//
// It uses github.com/francoispqt/gojay as JSON encoder.
//
// Basic usage:
// 	import "github.com/francoispqt/onelog/log"
//
//	log.Info("hello world !") // {"level":"info","message":"hello world !", "time":1494567715}
//
// You can create your own logger:
//	import "github.com/francoispqt/onelog
//
//	var logger = onelog.New(os.Stdout, onelog.ALL)
//
//	func main() {
//		logger.Info("hello world !") // {"level":"info","message":"hello world !"}
//	}
package onelog
