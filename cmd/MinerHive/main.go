package main

import (
	"log"

	"github.com/FedoraTipper/MinerHive/internal/config"
	hiveModule "github.com/FedoraTipper/MinerHive/internal/hive"
	configModels "github.com/FedoraTipper/MinerHive/internal/models/config"
)

var (
	configName  = "config.yml"
	configPaths = []string{
		".",
		"$HOME/minerhive/",
		"$HOME/.config/minerhive/",
	}
)

var defaultConfigValues = map[string]interface{}{
	"Port":         8080,
	"LoggingLevel": "info",
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
