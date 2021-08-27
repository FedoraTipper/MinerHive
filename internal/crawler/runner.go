package crawler

import (
	"fmt"
	"log"
	"time"

	"github.com/FedoraTipper/AntHive/internal/crawler/rpc"
	"github.com/FedoraTipper/AntHive/internal/logger"
	"github.com/FedoraTipper/AntHive/internal/models/config"
	"github.com/FedoraTipper/AntHive/internal/stasher"
	"github.com/FedoraTipper/AntHive/internal/transformer"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

type CrawlerRunner struct {
	stasher       *stasher.Stasher
	CrawlerConfig config.CrawlerConfig
}

func (cr *CrawlerRunner) StartWork() {
	cr.initLogger()

	s := gocron.NewScheduler(time.UTC)

	stash := &stasher.Stasher{}
	redisConfig := cr.CrawlerConfig.Redis

	err := stash.NewRedisClient(redisConfig.Host, redisConfig.Port, redisConfig.Username, redisConfig.Password, redisConfig.SelectedDatabase)

	if err != nil {
		zap.S().Fatalw("Error creating new Redis client", "Error", err)
	}

	zap.S().Infof("Redis client established to %s:%d", redisConfig.Host, redisConfig.Port)

	cr.stasher = stash

	for _, miner := range cr.CrawlerConfig.Miners {
		job, err := s.Every(fmt.Sprintf("%ds", cr.CrawlerConfig.CrawlInterval)).Do(cr.collect, miner)

		if err != nil {
			zap.S().Errorw("Error completing job for miner", "Miner", miner.MinerName, "Error", err)
			continue
		}
		// Set Scheduler in a singleton mode to avoid job collision
		job.SingletonMode()
	}

	s.StartBlocking()
}

func (cr *CrawlerRunner) initLogger() {
	err := logger.InitGlobalLogger(cr.CrawlerConfig.LoggingFile, cr.CrawlerConfig.LoggingLevel)

	if err != nil {
		log.Fatalf("Unable to configure logger. Error: %v", err)
	}
}

func (cr *CrawlerRunner) collect(miner config.MinerConfig) {
	zap.S().Infof("Starting new job for miner %s (%s)", miner.MinerName, miner.GetURL())

	rpcClient, err := rpc.GetRPCClient(miner.Model)

	if err != nil {
		zap.S().Errorw("Error getting RPC client for miner model", "Error", err)
		return
	}

	statsBytes, err := rpcClient.GetStats(miner.GetURL())

	if err != nil {
		zap.S().Errorw("Error getting stats for miner", "Miner", miner.MinerName, "Error", err)
		return
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
