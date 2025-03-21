package redis

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type RedisConfig struct {
	Address      string
	DB           int
	QueueName    string
	JobStatusTTL time.Duration
}

var cfg = initRedis()

func initRedis() RedisConfig {
	return RedisConfig{
		Address:      getEnv("REDIS_ADDRESS", "127.0.0.1:6379"),
		DB:           toInt(getEnv("REDIS_DB", "0")),
		QueueName:    getEnv("REDIS_QUEUE_NAME", "jobs"),
		JobStatusTTL: time.Duration(toInt(getEnv("REDIS_JOB_STATUS_TTL_SECONDS", "604800"))) * time.Second,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func toInt(value string) int {
	num, err := strconv.Atoi(value)
	if err != nil {
		log.Error().Msg("Redis: unable to convert env to int type")
		return 0
	}
	return num
}
