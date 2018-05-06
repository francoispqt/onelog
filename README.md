# onelog

onelog is a dead simple but very efficient JSON logger. 
It is one of the fastest JSON logger out there and the fastest when logging extra fields and keeps low allocation. 

It gives more control over log levels enabled by using bitwise operation for setting levels on a logger.

It is also modular as you can add a custom hook, define level and message keys.

Go 1.9 is required as it uses a type alias over gojay.Encoder.

It is named onelog because when logging just a message it does a single allocation (reference to zerolog), and also because it sounds like `One Love` song from Bob Marley :)

## Get Started

```bash 
go get github.com/francoispqt/onelog
```

Basic usage:

```go
import "github.com/francoispqt/onelog"

func main() {
    // create a new Logger
    // first argument is an io.Writer
    // second argument is the level, which is an integer
    logger := onelog.NewLogger(
        os.Stdout, 
        onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL,
    )
    logger.Info("hello world !") // {"level":"info","message":"hello world"}
}
```

## Levels

Levels are ints mapped to a string. The logger will check if level is enabled with a very efficient bitwise &, which makes onelog the fastest when running disabled logging with 0 allocs and merely 1ns/op. See benchmarks. 

When creating a logger you must use the `|` operator with different levels to toggle bytes. 

Example if you want levels INFO and WARN:
```go
logger := onelog.NewLogger(
    os.Stdout, 
    onelog.INFO|onelog.WARN,
)
```

This allows you to have a logger with different levels, for example you can do: 
```go
var logger *onelog.Logger

func init() {
    // if we are in debug mode, enable DEBUG lvl
    if os.Getenv("DEBUG") != "" {
        logger = onelog.NewLogger(
            os.Stdout, 
            onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL,
        )
        return
    }
    logger = onelog.NewLogger(
        os.Stdout, 
        onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL,
    )
}
```

Available levels:
- onelog.DEBUG
- onelog.INFO
- onelog.WARN
- onelog.ERROR
- onelog.FATAL

You can change their textual values by doing: 
```go
onelog.Levels[onelog.INFO] = "FOOBAR"
```

## Hook

You can define a hook which will be run for every log message. Example:
```go 
logger := onelog.NewLogger(
    os.Stdout, 
    onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL,
)
logger.Hook(func(enc *Encoder) {
    enc.AddStringKey("time", time.Now().Format(time.RFC3339))
})
logger.Info("hello world !") // {"level":"info","message":"hello world","time":"2018-05-06T02:21:01+08:00"}
```

## Logging

### Without extra fields
Logging without extra fields is easy as:
```go 
logger := onelog.NewLogger(
    os.Stdout, 
    onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL,
)
logger.Debug("i'm not sure what's going on") // {"level":"debug","message":"i'm not sure what's going on"}
logger.Info("breaking news !") // {"level":"info","message":"breaking news !"}
logger.Warn("beware !") // {"level":"warn","message":"beware !"}
logger.Error("my printer is on fire") // {"level":"error","message":"my printer is on fire"}
logger.Fatal("oh my...") // {"level":"fatal","message":"oh my..."}
```

### With extra fields 
Logging with extra fields is quite simple, specially if you have used gojay:
```go 
logger := onelog.NewLogger(
    os.Stdout, 
    onelog.DEBUG|onelog.INFO|onelog.WARN|onelog.ERROR|onelog.FATAL,
)

logger.DebugWithFields("i'm not sure what's going on", func(enc *onelog.Encoder) {
    enc.AddStringKey("userID", "123455")
    enc.AddObjectKey("user", onelog.ObjectFunc(func(enc *onelog.Encoder){
        enc.AddStringKey("name", "somename")
    }))
}) // {"level":"debug","message":"i'm not sure what's going on","userID":"123456","user":{"name":"somename"}}

logger.Info("breaking news !", func(enc *onelog.Encoder) {
    enc.AddStringKey("userID", "123455")
}) // {"level":"info","message":"breaking news !","userID":"123456"}

logger.Warn("beware !", func(enc *onelog.Encoder) {
    enc.AddStringKey("userID", "123455")
}) // {"level":"warn","message":"beware !","userID":"123456"}

logger.Error("my printer is on fire", func(enc *onelog.Encoder) {
    enc.AddStringKey("userID", "123455")
}) // {"level":"error","message":"my printer is on fire","userID":"123456"}

logger.Fatal("oh my...", func(enc *onelog.Encoder) {
    enc.AddStringKey("userID", "123455")
}) // {"level":"fatal","message":"oh my...","userID":"123456"}
```

# Benchmarks

The benchmark data presented here is the one from Uber's benchmark suite where we added onelog. 

A pull request will be submitted to Zap to integrate onelog in the benchmarks.

## Disabled Logging
|             | ns/op | bytes/op     | allocs/op |
|-------------|-------|--------------|-----------|
| Zap         | 7.70  | 0            | 0         |
| zerolog     | 2.67  | 0            | 0         |
| logrus      | 12.1  | 16           | 1         |
| onelog      | 0.88  | 0            | 0         |

## Disabled with fields
|             | ns/op | bytes/op     | allocs/op |
|-------------|-------|--------------|-----------|
| Zap         | 208   | 768          | 5         |
| zerolog     | 68.7  | 128          | 4         |
| logrus      | 721   | 1493         | 12        |
| onelog      | 1.31  | 0            | 0         |

## Logging basic message
|             | ns/op | bytes/op     | allocs/op |
|-------------|-------|--------------|-----------|
| Zap         | 203   | 0            | 0         |
| zerolog     | 154   | 0            | 0         |
| logrus      | 1256  | 1554         | 24        |
| onelog      | 329   | 0            | 0         |

## Logging message with extra fields
|             | ns/op | bytes/op     | allocs/op |
|-------------|-------|--------------|-----------|
| Zap         | 1764  | 770          | 5         |
| zerolog     | 9869  | 8515         | 84        |
| logrus      | 13211 | 13584        | 129       |
| onelog      | 1079  | 176          | 4         |