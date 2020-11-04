package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfigurations struct {
	Host     string
	GrpcPort string
	HttpPort string
}

var Server ServerConfigurations

func init() {
	err := viper.UnmarshalKey("server", &Server)
	if err != nil {
		log.Fatalf("Unable to decode into config struct, %v", err)
	}
}
