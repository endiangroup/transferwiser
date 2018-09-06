package main

import (
	"github.com/endiangroup/transferwiser/core"
	"github.com/endiangroup/transferwiser/web"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	config := core.GetConfig()

	transferwiseAPI := core.NewTransferwiseAPI(config.TwHost, config.TwAPIToken)

	webServer := web.NewServer(logger, transferwiseAPI)
	err = webServer.Run(config.Env, config.Port, config.LetsEncryptPort)
	if err != nil {
		logger.Error("error running web.", zap.Error(err))
	}
}
