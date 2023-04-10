package config

import (
	"ChatDanBackend/common"
	"github.com/caarlos0/env/v8"
	"go.uber.org/zap"
)

var Config struct {
	Debug    bool   `env:"DEBUG" envDefault:"false"`
	DbUrl    string `env:"DB_URL,required"`
	RedisUrl string `env:"REDIS_URL,required"`
}

func init() {
	if err := env.Parse(&Config); err != nil {
		panic(err)
	}
	common.Logger.Info("", zap.Any("config", Config))
}
