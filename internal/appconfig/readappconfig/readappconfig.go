package readappconfig

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
