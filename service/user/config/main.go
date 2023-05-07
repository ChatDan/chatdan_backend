package config

var CustomConfig struct {
	ApisixUrl      string `env:"APISIX_URL,required"`
	ApisixAdminKey string `env:"APISIX_ADMIN_KEY,required"`
	Hostname       string `env:"HOSTNAME,required"`
}
