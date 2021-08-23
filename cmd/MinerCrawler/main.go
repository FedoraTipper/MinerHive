package main

import (
	"github.com/FedoraTipper/AntHive/internal/config"
	"github.com/FedoraTipper/AntHive/internal/crawler"
	configModels "github.com/FedoraTipper/AntHive/internal/models/config"
	"log"
	//"github.com/uber-go/zap"
)

var (
	configName  = "config.yml"
	configPaths = []string{
		".",
		"./configs/minercrawler/",
		"./config/minercrawler/",
		"$HOME/.config/minercrawler/",
	}
)

var defaultConfigValues = map[string]interface{}{
	"CrawlInterval":  1, // 1 sec
	"LoggingEnabled": true,
	"LoggingLevel":   "info",
	"Salt":           "j8%mEbbkf2#PtSwxLnN4cSc&p%5SrJPviRDwGdrgx%STzW%P82s4j^e2PBtHvJ@J4%",
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
		log.Println("Error in config.yml")
		log.Fatalf("%v", err)
	}

	errs := crawlerConfig.Validate()

	for _, err := range errs {
		log.Printf("%v\n", err)
	}

	runner := crawler.CrawlerRunner{
		CrawlerConfig: crawlerConfig,
	}

	runner.StartWork()

}
