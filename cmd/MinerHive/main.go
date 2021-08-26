package main

import (
	"log"

	"github.com/FedoraTipper/AntHive/internal/config"
	hiveModule "github.com/FedoraTipper/AntHive/internal/hive"
	configModels "github.com/FedoraTipper/AntHive/internal/models/config"
)

var (
	configName  = "config.yml"
	configPaths = []string{
		".",
		"./configs/minerhive/",
		"./config/minerhive/",
		"$HOME/.config/minerhive/",
	}
)

var defaultConfigValues = map[string]interface{}{
	"Port":           8080,
	"LoggingEnabled": true,
	"LoggingLevel":   "info",
	"Redis": configModels.RedisConfig{
		SelectedDatabase: 0,
	},
}

func main() {
	viper := config.GenerateConfigReader(defaultConfigValues, configName, configPaths)

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error when reading in config.yml")
		log.Fatalf("%v", err)
	}

	var hiveConfig configModels.HiveConfig

	if err := viper.UnmarshalExact(&hiveConfig); err != nil {
		log.Println("Error in config.yml")
		log.Fatalf("%v", err)
	}

	errs := hiveConfig.Validate()

	for _, err := range errs {
		log.Printf("%v\n", err)
	}

	hive := hiveModule.Hive{
		Config: hiveConfig,
	}

	hive.Serve()
}
