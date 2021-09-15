package crawler

import (
	"fmt"
	"log"
	"time"

	"github.com/FedoraTipper/MinerHive/internal/crawler/rpc"
	"github.com/FedoraTipper/MinerHive/internal/logger"
	"github.com/FedoraTipper/MinerHive/internal/models/config"
	"github.com/FedoraTipper/MinerHive/internal/transformer"
	"github.com/FedoraTipper/MinerHive/pkg/constants"
	"github.com/FedoraTipper/MinerHive/pkg/models"
	"github.com/FedoraTipper/MinerHive/pkg/redis"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

type CrawlerRunner struct {
	wrappedRedisClient *redis.RedisClient
	CrawlerConfig      config.CrawlerConfig
}

func (cr *CrawlerRunner) StartWork() {
	cr.initLogger()

	s := gocron.NewScheduler(time.UTC)

	newWrappedRedisClient := &redis.RedisClient{}
	redisConfig := cr.CrawlerConfig.Redis

	err := newWrappedRedisClient.NewRedisClient(redisConfig.GetAddress(), redisConfig.Username, redisConfig.Password, redisConfig.SelectedDatabase)

	if err != nil {
		zap.S().Fatalw("Fatal error creating new Redis client", "Error", err)
	}

	zap.S().Infof("Redis client established to %s:%d", redisConfig.Host, redisConfig.Port)

	cr.wrappedRedisClient = newWrappedRedisClient

	for _, miner := range cr.CrawlerConfig.Miners {
		zap.S().Infof("Scheduling a job every %d seconds for %s", cr.CrawlerConfig.CrawlInterval, miner.MinerName)
		job, err := s.Every(fmt.Sprintf("%ds", cr.CrawlerConfig.CrawlInterval)).Do(cr.collect, miner)

		if err != nil {
			zap.S().Errorw("Error creating a scheduled job for miner", "Miner", miner.MinerName, "Error", err)
			continue
		}
		// Set Scheduler in a singleton mode to avoid job collision
		job.SingletonMode()
	}

	zap.S().Infof("Starting schedule with %d jobs", len(s.Jobs()))
	s.StartBlocking()
}

func (cr *CrawlerRunner) initLogger() {
	err := logger.InitGlobalLogger(cr.CrawlerConfig.LoggingFile, cr.CrawlerConfig.LoggingLevel)

	if err != nil {
		log.Fatalf("Unable to configure logger. Error: %v", err)
	}
}

func (cr *CrawlerRunner) collect(miner config.MinerConfig) {
	minerObj := &models.MinerStats{
		MinerName:   miner.MinerName,
		CrawlerName: cr.CrawlerConfig.CrawlerName,
		Status:      constants.MinerStatusUnknown,
	}

	zap.S().Infof("Starting new job for miner %s (%s)", miner.MinerName, miner.GetAddress())

	defer cr.stashMinerStats(&minerObj)

	rpcClient, err := rpc.GetRPCClient(miner.Model)

	if err != nil {
		zap.S().Errorw("Error getting RPC client for miner model", "Error", err)
		return
	}
	zap.S().Debugf("Successfully got RPC Client for %s", miner.Model)

	zap.S().Infof("Making call to RPC to get stats from %s", miner.GetAddress())
	statsBytes, err := rpcClient.GetStats(miner.GetAddress())

	if err != nil {
		zap.S().Errorw("Error getting stats for miner", "Miner", miner.MinerName, "URL", miner.GetAddress(), "Error", err)
		return
	}
	zap.S().Infof("Successfully got stats from RPC API - %s", miner.GetAddress())

	t, err := transformer.GetTransformer(miner.Model)

	if err != nil {
		zap.S().Errorw("Error getting stats transformer", "Miner", miner.MinerName, "Error", err)
		return
	}
	zap.S().Debugf("Successfully got stats transformer for %s", miner.Model)

	minerObj, err = t.ConvertStatsPayloadToMinerStats(miner.MinerName, cr.CrawlerConfig.CrawlerName, statsBytes)

	if err != nil {
		zap.S().Errorw("Error converting stats payload to miner stats model", "Miner", miner.MinerName, "JSON", string(statsBytes), "Error", err)
		return
	}
	zap.S().Debug("Successfully converted RPC stats model to MinerStats model", "Miner", miner.MinerName)
}

func (cr *CrawlerRunner) stashMinerStats(minerStatPtrAddress **models.MinerStats) {
	// dereference address to minerstats obj
	minerStats := *minerStatPtrAddress

	expiration, err := time.ParseDuration(fmt.Sprintf("%ds", cr.CrawlerConfig.CrawlInterval))

	if err != nil {
		zap.S().Fatalw("Unable to parse duration of Crawl Interval", "Error", err)
	}

	expiration *= 2

	zap.S().Infow("Stashing MinerStats interface in RedisDB with expiration", "Key", minerStats.MinerName, "Expiration", expiration.String())
	err = cr.wrappedRedisClient.StashInterface(minerStats.MinerName, minerStats, expiration)

	if err != nil {
		zap.S().Errorw("Error inserting miner stats model into RedisDB", "Miner", minerStats.MinerName, "Error", err)
	}
}
