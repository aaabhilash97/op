package logger

import (
	"log"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (

	// Log is global logger
	Log *zap.Logger

	// accessLogonceInit guarantee initialize logger only once
	onceInit sync.Once
)

// Info - log to info level
var Info func(msg string, fields ...zap.Field)

// Error - log to error level
var Error func(msg string, fields ...zap.Field)

// Debug - log to debug level
var Debug func(msg string, fields ...zap.Field)

// Warn  - log to warn level
var Warn func(msg string, fields ...zap.Field)

// Fatal  - log to Fatal level
var Fatal func(msg string, fields ...zap.Field)

// Initialize default logger with stdout and debug level
func init() {
	logger, err := CreateLogger(
		LogLevels["debug"],
		"stdout",
	)
	if err != nil {
		log.Fatal("Failed initialize logger", err)
	}
	Log = logger
	Info = logger.Info
	Error = logger.Error
	Debug = logger.Debug
	Warn = logger.Warn
	Fatal = logger.Fatal
}

func Init(level, output string) {
	onceInit.Do(func() {
		logger, err := CreateLogger(
			LogLevels[level],
			output,
		)
		if err != nil {
			log.Fatal("Failed initialize logger", err)
		}
		Log = logger
		Info = logger.Info
		Error = logger.Error
		Debug = logger.Debug
		Warn = logger.Warn
		Fatal = logger.Fatal
	})
}

// Init initializes log by input parameters
// lvl - global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
func CreateLogger(lvl int, logOutputPath string) (*zap.Logger, error) {

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	loggerConfig.DisableStacktrace = true
	loggerConfig.OutputPaths = strings.Split(logOutputPath, ",")
	loggerConfig.Level = zap.NewAtomicLevelAt(
		zapcore.Level(lvl),
	)

	logger, err := loggerConfig.Build()
	return logger, err

}

// LogLevels is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
var LogLevels = map[string]int{
	"debug":  -1,
	"info":   0,
	"warn":   1,
	"error":  2,
	"dpanic": 3,
	"panic":  4,
	"fatal":  5,
}
