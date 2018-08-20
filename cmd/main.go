package main

import (
	"github.com/endiangroup/transferwiser/core"
	"github.com/endiangroup/transferwiser/keyvalue"
	"github.com/endiangroup/transferwiser/web"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	config := core.GetConfig()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisKV := keyvalue.NewRedisKeyValue(redisClient)
	authKeyStore := keyvalue.NewValue(redisKV, "authentication_token")

	transferwiseAPI := core.NewTransferwiseAPI()

	transferwiseService := core.NewTransferwiseService(transferwiseAPI, config.TWClientID, config.TWHost, config.TWLoginRedirect, authKeyStore)

	webServer := web.NewServer(logger, transferwiseService)
	err = webServer.Run(config.Port)
	if err != nil {
		logger.Error("error running web.", zap.Error(err))
	}
}
