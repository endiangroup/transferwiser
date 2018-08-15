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

	transferwiseService := core.NewTransferwiseService(config.TWClientID, config.TWHost, config.TWLoginRedirect)

	webServer := web.NewServer(logger, transferwiseService)
	err = webServer.Run(config.Port)
	if err != nil {
		logger.Error("error running web.", zap.Error(err))
	}
}
