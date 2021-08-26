package hive

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FedoraTipper/AntHive/internal/models/config"
	stasher2 "github.com/FedoraTipper/AntHive/internal/stasher"
	"github.com/gin-gonic/gin"
)

type Hive struct {
	Config config.HiveConfig
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

func (h *Hive) Serve() {
	router := gin.Default()
	router.Use(h.AuthorizationChallengeMiddleware)

	router.GET(fmt.Sprintf("/:%s/stats", minerIdParam), h.getMinerStats)
	router.GET(fmt.Sprintf("/miners", minerIdParam), h.getMinerStats)

	router.Run(fmt.Sprintf(":%d", h.Config.Port))
}

func (h *Hive) getMinerStats(c *gin.Context) {
	minerName := c.Param(minerIdParam)
	stasher := &stasher2.Stasher{}
	redisConfig := h.Config.Redis
	err := stasher.NewRedisClient(redisConfig.Host, redisConfig.Port, redisConfig.Username, redisConfig.Password, redisConfig.SelectedDatabase)

	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "An error occured")
		return
	}

	miner, err := stasher.GetInterface(minerName)

	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, "An error occured")
		return
	}

	if miner == "" {
		c.String(http.StatusNotFound, fmt.Sprintf("%s does not exist in key db", minerName))
		return
	}

	c.String(http.StatusOK, miner)
}
