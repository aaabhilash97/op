package appconfig

import (
	"log"

	_ "github.com/aaabhilash97/op/internal/appconfig/readappconfig"
	"github.com/aaabhilash97/op/pkg/config"
	"github.com/aaabhilash97/op/pkg/logger"
	"github.com/spf13/viper"
)

var Appconfig config.Config

func init() {
	err := viper.Unmarshal(Appconfig)
	if err != nil {
		log.Fatalf("Unable to decode into config struct, %v", err)
	}

	// Initialize logger, on any error fn will panic
	logger.Init(Appconfig.Log.LogLevel, Appconfig.Log.LogOutputPaths)
}
