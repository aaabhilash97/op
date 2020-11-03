package config

import (
	"log"

	_ "github.com/aaabhilash97/op/pkg/config/readconfig"
	"github.com/spf13/viper"
)

type Configurations struct {
	Server ServerConfigurations
}
type ServerConfigurations struct {
	Host     string
	GrpcPort string
	RestPort string
}

var Config Configurations

func init() {
	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("Unable to decode into config struct, %v", err)
	}
}
