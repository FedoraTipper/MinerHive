package crawler

import (
	"fmt"
	"log"
	"time"

	"github.com/FedoraTipper/AntHive/internal/crawler/rpc"
	"github.com/FedoraTipper/AntHive/internal/models/config"
	"github.com/FedoraTipper/AntHive/internal/stasher"
	"github.com/FedoraTipper/AntHive/internal/transformer"
	"github.com/go-co-op/gocron"
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
		job, err := s.Every(fmt.Sprintf("%ds", cr.CrawlerConfig.CrawlInterval)).Do(cr.collect, miner)

		if err != nil {
			log.Println(err)
			continue
		}
		// Set Scheduler in a singleton mode to avoid job collision
		job.SingletonMode()
	}

	s.StartBlocking()
}

func (cr *CrawlerRunner) collect(miner config.MinerConfig) {
	fmt.Printf("Starting job for miner %s\n", miner.MinerName)
	//url := http.FormURL(miner.Host, miner.Port)
	url := fmt.Sprintf("%s:%d", miner.Host, miner.Port)

	rpcClient, err := rpc.GetRPCClient(miner.Model)

	if err != nil {
		log.Fatalln(err)
	}

	statsBytes, err := rpcClient.GetStats(url)

	if err != nil {
		log.Fatalln(err)
	}

	t, err := transformer.GetTransformer(miner.Model)

	if err != nil {
		log.Fatalln(err)
	}

	minerObj, err := t.ConvertStatsPayloadToMiner(miner.MinerName, cr.CrawlerConfig.CrawlerName, statsBytes)

	if err != nil {
		log.Fatalln(err)
	}

	err = cr.stasher.StashInterface(minerObj)

	if err != nil {
		log.Fatal(err)
	}
}
