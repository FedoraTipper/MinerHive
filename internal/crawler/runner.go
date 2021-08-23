package crawler

import (
	"fmt"
	"github.com/FedoraTipper/AntHive/internal/constants"
	http2 "github.com/FedoraTipper/AntHive/internal/crawler/http"
	"github.com/FedoraTipper/AntHive/internal/models/config"
	"github.com/FedoraTipper/AntHive/internal/stasher"
	"github.com/FedoraTipper/AntHive/internal/transformer"
	"github.com/FedoraTipper/AntHive/pkg/http"
	"github.com/go-co-op/gocron"
	"io/ioutil"
	"log"
	"time"
)

type CrawlerRunner struct {
	stasher       *stasher.Stasher
	CrawlerConfig config.CrawlerConfig
}

func (cr *CrawlerRunner) StartWork() {
	s := gocron.NewScheduler(time.UTC)

	stash := &stasher.Stasher{}
	redisConfig := cr.CrawlerConfig.Redis

	err := stash.NewRedisClient(redisConfig.Host, redisConfig.Port, redisConfig.Username, redisConfig.Password, redisConfig.SelectedDatabase)

	if err != nil {
		log.Fatalln(err)
	}

	cr.stasher = stash

	for _, miner := range cr.CrawlerConfig.Miners {
		job, err := s.Every(fmt.Sprintf("%ds", cr.CrawlerConfig.CrawlInterval)).Do(cr.crawl, miner)

		if err != nil {
			log.Println(err)
			continue
		}
		// Set Scheduler in a singleton mode to avoid job collision
		job.SingletonMode()
	}

	s.StartBlocking()
}

func (cr *CrawlerRunner) crawl(miner config.MinerConfig) {
	fmt.Printf("Starting job for miner %s\n", miner.MinerName)
	url := http.FormURL(miner.Host, miner.Port)

	httpClient, err := http2.GetHTTPClient(miner.Model)

	if err != nil {
		log.Fatalln(err)
	}

	statsResp, err := httpClient.GetStatsResponse(url, miner.Username, miner.Password, cr.CrawlerConfig.Salt)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(statsResp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	statsResp.Body.Close()

	statsPayload, err := parseS19Stats(body)

	if err != nil {
		log.Fatalln(err)
	}

	t, err := transformer.GetTransformer(constants.X19)

	if err != nil {
		log.Fatalln(err)
	}

	minerObj, err := t.ConvertStatsPayloadToMiner(miner.MinerName, cr.CrawlerConfig.CrawlerName, statsPayload)

	if err != nil {
		log.Fatalln(err)
	}

	err = cr.stasher.StashInterface(minerObj)

	if err != nil {
		log.Fatal(err)
	}
}
