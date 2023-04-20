package configx

import (
	"ChatDanBackend/common"
	"github.com/caarlos0/env/v8"
	"go.uber.org/zap"
)

var CommonConfig struct {
	Debug    bool   `env:"DEBUG" envDefault:"false"`
	Mode     string `env:"MODE" envDefault:"dev"`
	DbUrl    string `env:"DB_URL,required"`
	RedisUrl string `env:"REDIS_URL,required"`
}

func InitConfig(CustomConfig any) {
	if err := env.Parse(&CommonConfig); err != nil {
		panic(err)
	}
	common.Logger.Info("", zap.Any("common_config", CommonConfig))

	if CustomConfig != nil {
		if err := env.Parse(&CustomConfig); err != nil {
			panic(err)
		}
		common.Logger.Info("", zap.Any("custom_config", CustomConfig))
	}
}
