package util

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"urlshortener/model"
)

func CheckCache(s string) (string, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()

	return rdb.Get(ctx, s).Result()
}

func WriteCache(kvPair model.KvPair) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()

	// Write with a TTL (expiration)
	jsonBytes, _ := json.Marshal(kvPair.Original)
	rdb.Set(ctx, kvPair.Shortened, jsonBytes, 5*time.Minute)
}
