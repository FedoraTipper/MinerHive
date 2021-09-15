package hive

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FedoraTipper/MinerHive/internal/logger"
	"github.com/FedoraTipper/MinerHive/internal/models/config"
	redisWrapper "github.com/FedoraTipper/MinerHive/pkg/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Hive struct {
	Config      config.HiveConfig
	redisClient *redisWrapper.RedisClient
}

const (
	apiKeyHeader = "api-key"
	minerIdParam = "minerName"
)

func (h *Hive) AuthorizationChallengeMiddleware(c *gin.Context) {
	apiKey := c.GetHeader(apiKeyHeader)

	if apiKey == h.Config.Token {
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
	c.String(http.StatusUnauthorized, "Request unauthorized")
}

// TODO: [Tidy up, Split into different files, Wrap responses (payloads and errors)]

func (h *Hive) Serve() {
	h.initLogger()

	redis := &redisWrapper.RedisClient{}
	redisConfig := h.Config.Redis
	err := redis.NewRedisClient(redisConfig.GetAddress(), redisConfig.Username, redisConfig.Password, redisConfig.SelectedDatabase)

	if err != nil {
		zap.S().Fatalw("Fatal error creating new Redis client", "Error", err)
	}

	h.redisClient = redis
	router := gin.Default()
	router.Use(h.AuthorizationChallengeMiddleware)

	router.GET(fmt.Sprintf("/miner/:%s/stats", minerIdParam), h.getMinerStats)
	router.GET(fmt.Sprintf("/miners"), h.getMinerKeys)

	router.Run(fmt.Sprintf(":%d", h.Config.Port))
}

func (h *Hive) initLogger() {
	err := logger.InitGlobalLogger(h.Config.LoggingFile, h.Config.LoggingLevel)

	if err != nil {
		log.Fatalf("Unable to configure logger. Error: %v", err)
	}
}

func (h *Hive) getMinerStats(c *gin.Context) {
	minerName := c.Param(minerIdParam)

	if len(minerName) == 0 {
		zap.S().Warnf("Missing %s in request parameters from %s", minerIdParam, c.Request.RemoteAddr)
		c.String(http.StatusBadRequest, "Malformed request. Miner name not set as parameter")
		return
	}

	miner, err := h.redisClient.GetInterface(minerName)

	if err != nil {
		zap.S().Errorw("Unable to get miner from RedisDB", "MinerName", minerName, "Error", err)
		c.String(http.StatusInternalServerError, "An error occured")
		return
	}

	if miner == "" {
		c.String(http.StatusNotFound, fmt.Sprintf("%s does not exist in key db", minerName))
		return
	}

	c.String(http.StatusOK, miner)
}

func (h *Hive) getMinerKeys(c *gin.Context) {
	zap.S().Info("Getting all keys from RedisDB")
	minerKeys, err := h.redisClient.GetKeys()

	if err != nil {
		zap.S().Errorw("Unable to get miner keys from RedisDB", "Error", err)
		c.String(http.StatusInternalServerError, "An error occured")
		return
	}

	if len(minerKeys) == 0 {
		zap.S().Warn("No keys returned from RedisDB")
		c.String(http.StatusNotFound, "No keys exist in RedisDB")
		return
	}

	c.JSON(http.StatusOK, minerKeys)
}
