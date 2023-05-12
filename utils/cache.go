package utils

import (
	"ChatDanBackend/config"
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/goccy/go-json"
	"github.com/juju/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var usingRedis = false
var RedisClient *redis.Client
var BigCacheClient *bigcache.BigCache

func InitCache() {
	if config.Config.RedisUrl != "" {
		RedisClient = redis.NewClient(&redis.Options{Addr: config.Config.RedisUrl})
		pong := RedisClient.Ping(context.Background())
		if pong.Err() != nil {
			panic(pong.Err())
		} else {
			Logger.Info("redis ping success")
		}
		usingRedis = true
	} else {
		Logger.Info("redis url not set, using bigcache")
		var err error
		BigCacheClient, err = bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
		if err != nil {
			panic(err)
		}
	}

}

func Get(key string, model any) (err error) {
	var value []byte
	if usingRedis {
		value, err = RedisClient.Get(context.Background(), key).Bytes()
	} else {
		value, err = BigCacheClient.Get(key)
	}
	if err != nil {
		return errors.Trace(err)
	}
	err = json.Unmarshal(value, model)
	return errors.Trace(err)
}

func Set(key string, model any) (err error) {
	var value []byte
	value, err = json.Marshal(model)
	if err != nil {
		return errors.Trace(err)
	}

	if usingRedis {
		return errors.Trace(RedisClient.Set(context.Background(), key, value, 0).Err())
	} else {
		return errors.Trace(BigCacheClient.Set(key, value))
	}
}

func Delete(key string) error {
	if usingRedis {
		return errors.Trace(RedisClient.Del(context.Background(), key).Err())
	} else {
		return errors.Trace(BigCacheClient.Delete(key))
	}
}
