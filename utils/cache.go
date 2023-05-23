package utils

import (
	"chatdan_backend/config"
	"context"
	"github.com/allegro/bigcache/v3"
	"github.com/goccy/go-json"
	"github.com/juju/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

var usingRedis = false
var RedisClient *redis.Client
var BigCacheClient *bigcache.BigCache

var ErrCacheMiss = errors.New("cache miss")

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
	if config.Config.Debug {
		Logger.Info("get cache", zap.String("key", key))
	}
	var value []byte
	if usingRedis {
		value, err = RedisClient.Get(context.Background(), key).Bytes()
	} else {
		value, err = BigCacheClient.Get(key)
	}
	if err != nil {
		if err == redis.Nil || err == bigcache.ErrEntryNotFound {
			return ErrCacheMiss
		}
		return
	}
	return json.Unmarshal(value, model)
}

func Set(key string, model any, expiration time.Duration) (err error) {
	if config.Config.Debug {
		defer Logger.Info("set cache", zap.String("key", key), zap.Error(err))
	}
	var value []byte
	value, err = json.Marshal(model)
	if err != nil {
		return errors.Trace(err)
	}

	if usingRedis {
		return errors.Trace(RedisClient.Set(context.Background(), key, value, expiration).Err())
	} else {
		return errors.Trace(BigCacheClient.Set(key, value))
	}
}

func Delete(key string) {
	if usingRedis {
		_ = RedisClient.Del(context.Background(), key)
	} else {
		_ = BigCacheClient.Delete(key)
	}
}

func DeleteInBatch(keys ...string) {
	if len(keys) == 0 {
		return
	}
	if config.Config.Debug {
		Logger.Debug("delete in batch", zap.Strings("keys", keys))
	}
	if usingRedis {
		_ = RedisClient.Del(context.Background(), keys...)
	} else {
		for _, key := range keys {
			_ = BigCacheClient.Delete(key)
		}
	}
}
