package core

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Port int `default:"8080"`

	TWHost          string `default:"sandbox.transferwise.tech"`
	TWLoginRedirect string `required:"true"`
	TWClientID      string `default:"endiangroup/transferwiser"`

	RedisAddr string `default:"localhost:6379"`
}

var configInstance *config
var configOnce sync.Once

func GetConfig() *config {
	configOnce.Do(func() {
		var cfg config
		err := envconfig.Process("transferwiser", &cfg)
		if err != nil {
			panic(err)
		}
		configInstance = &cfg
	})
	return configInstance
}
