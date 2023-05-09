package utils

import (
	"ChatDanBackend/config"
	"context"
	"github.com/goccy/go-json"
	"github.com/juju/errors"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitCache() {
	RedisClient = redis.NewClient(&redis.Options{Addr: config.Config.RedisUrl})
}

func Get(key string, model any) error {
	value, err := RedisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		return errors.Trace(err)
	}
	err = json.Unmarshal(value, model)
	return errors.Trace(err)
}

func Set(key string, model any) error {
	val, err := json.Marshal(model)
	if err != nil {
		return errors.Trace(err)
	}
	return errors.Trace(RedisClient.Set(context.Background(), key, val, 0).Err())
}

func Delete(key string) error {
	return errors.Trace(RedisClient.Del(context.Background(), key).Err())
}
