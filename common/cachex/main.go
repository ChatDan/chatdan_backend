package cachex

import (
	"ChatDanBackend/common/configx"
	"context"
	"github.com/goccy/go-json"
	"github.com/juju/errors"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitCache() {
	RedisClient = redis.NewClient(&redis.Options{Addr: configx.CommonConfig.RedisUrl})
}

func Get[T any](key string) (*T, error) {
	var model T
	val, err := RedisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, errors.Trace(err)
	}
	err = json.Unmarshal(val, &model)
	return &model, errors.Trace(err)
}

func Set[T any](key string, model T) error {
	val, err := json.Marshal(model)
	if err != nil {
		return errors.Trace(err)
	}
	return errors.Trace(RedisClient.Set(context.Background(), key, val, 0).Err())
}

func Delete(key string) error {
	return errors.Trace(RedisClient.Del(context.Background(), key).Err())
}
