package main

import (
	"log"

	"github.com/barantoraman/microgate/pkg/config"
)

func main() {
	var cfg config.AuthServiceConfigurations
	loader := config.GetLoader("dev")
	if loader == nil {
		log.Fatal("cannot get env loader")
	}
	if err := loader.GetConfigByKey("auth_service", &cfg); err != nil {
		log.Fatal("cannot get env loader")
	}
}
