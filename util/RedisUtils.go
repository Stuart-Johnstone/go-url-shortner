package util

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"urlshortener/model"
)

func CheckCache(s string) (string, error) {
	rdb := getRedisClient()

	ctx := context.Background()
	res, err := rdb.Get(ctx, s).Result()

	var original string
	json.Unmarshal([]byte(res), &original)

	rdb.Close()

	return original, err
}

func WriteCache(kvPair model.KvPair) {
	rdb := getRedisClient()

	ctx := context.Background()

	// Write with a TTL (expiration)
	jsonBytes, _ := json.Marshal(kvPair.Original)
	rdb.Set(ctx, kvPair.Shortened, jsonBytes, 5*time.Minute)

	rdb.Close()
}

func getRedisClient() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	return redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}
