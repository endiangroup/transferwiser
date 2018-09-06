package core

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Port            int    `default:"8080"`
	LetsEncryptPort int    `default:"8081"`
	Env             string `default:"dev"`
	TwHost          string `default:"api.sandbox.transferwise.tech"`
	TwProfile       string `required:"true"`
	TwAPIToken      string `required:"true"`
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
