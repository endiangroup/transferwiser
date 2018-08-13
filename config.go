package transferwiser

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Port int `default:"8080"`
}

var configInstance *config
var configOnce sync.Once

func GetConfig() *config {
	configOnce.Do(func() {
		var cfg config
		err := envconfig.Process("transferwise", &cfg)
		if err != nil {
			panic(err)
		}
		configInstance = &cfg
	})
	return configInstance
}
