package main

import (
	"log"

	"github.com/FedoraTipper/AntHive/internal/config"
	"github.com/FedoraTipper/AntHive/internal/crawler"
	configModels "github.com/FedoraTipper/AntHive/internal/models/config"
)

var (
	configName  = "config.yml"
	configPaths = []string{
		".",
		"$HOME/minercrawler/",
		"$HOME/.config/minercrawler/",
	}
)

var defaultConfigValues = map[string]interface{}{
	"CrawlInterval": 1, // 1 sec
	"LoggingLevel":  "info",
	"LoggingFile":   "",
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

	var crawlerConfig configModels.CrawlerConfig

	if err := viper.UnmarshalExact(&crawlerConfig); err != nil {
		log.Println("Error in parsing config.yml")
		log.Fatalf("%v", err)
	}

	errs := crawlerConfig.Validate()

	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("%v\n", err)
		}

		log.Fatalln("End of errors")
	}

	runner := crawler.CrawlerRunner{
		CrawlerConfig: crawlerConfig,
	}

	runner.StartWork()
}
