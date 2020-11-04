package logger

import (
	"fmt"
	"log"
	"strings"
	"sync"

	_ "github.com/aaabhilash97/op/pkg/config/readconfig"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	LogLevel             string
	LogOutputPaths       string
	AccessLogLevel       string
	AccessLogOutputPaths string
}

var (

	// Access Log
	AccessLog *zap.Logger

	// Log is global logger
	Log *zap.Logger

	// accessLogonceInit guarantee initialize logger only once
	accessLogonceInit sync.Once
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

var Config LogConfig

func init() {

	viper.SetDefault("logger.LogLevel", "debug")
	viper.SetDefault("logger.LogOutputPaths", "stdout")
	viper.SetDefault("logger.AccessLogLevel", "error")
	viper.SetDefault("logger.AccessLogOutputPaths", "stdout")

	err := viper.UnmarshalKey("logger", &Config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	fmt.Println(Config)
	logger := initLogger(
		LogLevels[Config.LogLevel],
		Config.LogOutputPaths,
	)
	Log = logger
	Info = logger.Info
	Error = logger.Error
	Debug = logger.Debug
	Warn = logger.Warn
	Fatal = logger.Fatal
}

// Init initializes log by input parameters
// lvl - global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
func initLogger(lvl int, logOutputPath string) *zap.Logger {

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
	return logger

}

func InitAccessLogger(lvl int, logOutputPath string) {
	accessLogonceInit.Do(func() {
		logger := initLogger(LogLevels[Config.AccessLogLevel], Config.AccessLogOutputPaths)
		AccessLog = logger
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
