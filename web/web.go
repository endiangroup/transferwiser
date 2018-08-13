package web

import (
	"fmt"

	"github.com/endiangroup/transferwiser"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger) error {
	config := transferwiser.GetConfig()

	e := echo.New()

	addr := fmt.Sprintf(":%v", config.Port)
	logger.Info("Listening", zap.String("addr", addr))
	return e.Start(addr)
}
