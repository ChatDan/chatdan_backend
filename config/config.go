package config

import (
	"github.com/caarlos0/env/v8"
)

var Config struct {
	Debug             bool   `env:"DEBUG" envDefault:"false"`
	Mode              string `env:"MODE" envDefault:"dev"`
	DbType            string `env:"DB_TYPE" envDefault:"sqlite"`
	DbUrl             string `env:"DB_URL"`
	RedisUrl          string `env:"REDIS_URL"`
	AppName           string `env:"APP_NAME" envDefault:"ChatDan"`
	Hostname          string `env:"HOSTNAME" envDefault:"localhost"`
	Standalone        bool   `env:"STANDALONE" envDefault:"false"` // if true, go without gateway
	GatewayType       string `env:"GATEWAY_TYPE" envDefault:"apisix"`
	ApisixUrl         string `env:"APISIX_URL"`
	ApisixAdminKey    string `env:"APISIX_ADMIN_KEY"`
	MeilisearchUrl    string `env:"MEILISEARCH_URL"`
	MeilisearchApiKey string `env:"MEILISEARCH_API_KEY"`
}

func InitConfig() {
	var err error
	if err = env.Parse(&Config); err != nil {
		panic(err)
	}

	if !Config.Standalone {
		if Config.GatewayType == "apisix" {
			if Config.ApisixUrl == "" {
				panic("APISIX_URL is required")
			}
			if Config.ApisixAdminKey == "" {
				panic("APISIX_ADMIN_KEY is required")
			}
		} else {
			panic("unknown gateway type")
		}
	}
}
