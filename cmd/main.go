package main

import (
	"github.com/endiangroup/transferwiser/web"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	err = web.Run(logger)
	if err != nil {
		logger.Error("error running web.", zap.Error(err))
	}
}
