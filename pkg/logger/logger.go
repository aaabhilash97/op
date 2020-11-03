package logger

import (
	"log"
	"strings"
	"sync"

	_ "github.com/aaabhilash97/op/pkg/config/readconfig"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logConfig struct {
	LogLevel             string
	LogOutputPaths       string
	AccessLogLevel       string
	AccessLogOutputPaths string
}

var (
	// AccessLog access logger
	AccessLog *zap.Logger
	// Log is global logger
	Log *zap.Logger

	// onceInit guarantee initialize logger only once
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

func init() {
	var config logConfig

	err := viper.UnmarshalKey("logger", &config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	Init(
		LogLevels[config.LogLevel],
		config.LogOutputPaths,
		LogLevels[config.AccessLogLevel],
		config.AccessLogOutputPaths,
	)
}

// Init initializes log by input parameters
// lvl - global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
func Init(lvl int, logOutputPath string, accessLogLevel int, accessLogOutputPath string) {

	onceInit.Do(func() {
		{
			if logOutputPath == "" {
				logOutputPath = "stdout"
			}
			loggerConfig := zap.NewProductionConfig()
			loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			loggerConfig.DisableStacktrace = true
			loggerConfig.OutputPaths = strings.Split(logOutputPath, ",")
			loggerConfig.Level = zap.NewAtomicLevelAt(
				zapcore.Level(lvl),
			)

			logger, err := loggerConfig.Build()
			if err != nil {
				log.Fatal(err)
			}
			Log = logger
			Info = logger.Info
			Error = logger.Error
			Debug = logger.Debug
			Warn = logger.Warn
			Fatal = logger.Fatal
		}

		{
			if accessLogOutputPath == "" {
				accessLogOutputPath = "stdout"
			}
			loggerConfig := zap.NewProductionConfig()
			loggerConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
			loggerConfig.DisableCaller = true
			loggerConfig.DisableStacktrace = true
			loggerConfig.OutputPaths = strings.Split(
				accessLogOutputPath, ",")

			loggerConfig.Level = zap.NewAtomicLevelAt(
				zapcore.Level(accessLogLevel),
			)

			logger, err := loggerConfig.Build()
			if err != nil {
				log.Fatal(err)
			}
			AccessLog = logger
		}
	})

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
