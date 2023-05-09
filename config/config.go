package config

import (
	"github.com/caarlos0/env/v8"
)

var Config struct {
	Debug          bool   `env:"DEBUG" envDefault:"false"`
	DbUrl          string `env:"DB_URL,required"`
	RedisUrl       string `env:"REDIS_URL,required"`
	AppName        string `env:"APP_NAME" envDefault:"ChatDan"`
	Hostname       string `env:"HOSTNAME" envDefault:"localhost"`
	ApisixUrl      string `env:"APISIX_URL,required"`
	ApisixAdminKey string `env:"APISIX_ADMIN_KEY,required"`
}

func InitConfig() {
	var err error
	if err = env.Parse(&Config); err != nil {
		panic(err)
	}
}
