package main

// Temp script to feed MinerStats to InfluxDB bucket

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/FedoraTipper/AntHive/pkg/models"
	"github.com/go-redis/redis/v8"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/spf13/pflag"
)

// TODO: Move to own project when complete
func main() {
	var token, address, org, bucket, redisAddress, redisUsername, redisPassword, redisSelectedDatabase, minerName string

	pflag.StringVar(&minerName, "minername", "", "The miner name you would like to report on")
	pflag.StringVar(&token, "token", "", "Token for influxdb")
	pflag.StringVar(&address, "address", "", "Host and port to influxdb. E.g. influxdb.example.com:443")
	pflag.StringVar(&org, "org", "", "Influxdb organisation name")
	pflag.StringVar(&bucket, "bucket", "Miners", "Bucket to drop miner data into (Default: Miners)")
	pflag.StringVar(&redisAddress, "raddress", "", "Host and port to Redis DB. E.g. redis.example.com:6379")
	pflag.StringVar(&redisUsername, "rusername", "", "Username for RedisDB. Don't add flag if none is set")
	pflag.StringVar(&redisPassword, "rpassword", "", "Password for RedisDB. Don't add flag if none is set")
	pflag.StringVar(&redisSelectedDatabase, "rselecteddatabase", "0", "Redis selected DB. Don't add flag if none is set (Default: 1)")
	pflag.Parse()

	redisSelectedDatabaseInt, err := strconv.Atoi(redisSelectedDatabase)

	if err != nil {
		redisSelectedDatabaseInt = 0
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Username: redisUsername,
		Password: redisPassword,
		DB:       redisSelectedDatabaseInt,
	})

	err = redisTestConnection(redisClient)

	if err != nil {
		log.Fatal(err)
	}

	minerJson, err := getInterface(minerName, redisClient)

	if err != nil {
		log.Fatal(err)
	}

	if minerJson == "" {
		log.Fatal("No miner information in redis store")
	}

	miner := models.Miner{}

	err = miner.UnmarshalBinary([]byte(minerJson))

	if err != nil {
		log.Fatal(err)
	}

	// get non-blocking write client
	client := influxdb2.NewClient(address, token)

	// always close client at the end
	defer client.Close()
	writeAPI := client.WriteAPI(org, bucket)

	// write line protocol
	fanPoints := influxdb2.NewPointWithMeasurement("fans").
		AddTag("miner", minerName)

	for _, fan := range miner.Fans {
		fanPoints.AddField(fmt.Sprintf("fan%d", fan.FanNumber), fan.RPM)
	}

	writeAPI.WritePoint(fanPoints)

	var totalHashRate float64

	for _, board := range miner.HashBoards {
		pcbTempPoints := influxdb2.NewPointWithMeasurement("pcb-temperature").
			AddTag("miner", minerName).
			AddTag("board", fmt.Sprintf("%d", board.BoardNumber)).
			AddField("inlet", math.Max(float64(board.ChipTemperature[0]), float64(board.ChipTemperature[1]))).
			AddField("outlet", math.Max(float64(board.PCBTemperature[2]), float64(board.PCBTemperature[3])))

		chipTempPoints := influxdb2.NewPointWithMeasurement("chip-temperature").
			AddTag("miner", minerName).
			AddTag("board", fmt.Sprintf("%d", board.BoardNumber)).
			AddField("inlet", math.Max(float64(board.ChipTemperature[0]), float64(board.ChipTemperature[1]))).
			AddField("outlet", math.Max(float64(board.ChipTemperature[2]), float64(board.ChipTemperature[3])))

		writeAPI.WritePoint(chipTempPoints)
		writeAPI.WritePoint(pcbTempPoints)

		totalHashRate += board.CurrentHashRate
	}

	hashratePoint := influxdb2.NewPointWithMeasurement("hashrate").AddTag("miner", minerName).AddField("total", totalHashRate)
	uptimePoint := influxdb2.NewPointWithMeasurement("uptime").AddTag("miner", minerName).AddField("uptime", miner.Uptime)

	writeAPI.WritePoint(hashratePoint)
	writeAPI.WritePoint(uptimePoint)

	// Flush writes
	writeAPI.Flush()
}

func redisTestConnection(client *redis.Client) error {
	ctx := context.Background()

	return client.Ping(ctx).Err()
}

func getInterface(key string, redisClient *redis.Client) (string, error) {
	ctx := context.Background()

	i, err := redisClient.Get(ctx, key).Result()

	if err == redis.Nil {
		err = nil
	} else if err != nil {
		return "", err
	}

	return i, nil
}
