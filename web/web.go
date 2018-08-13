package web

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/endiangroup/transferwiser/core"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

var redirectUrl string
var redirectUrlOnce sync.Once

func Run(logger *zap.Logger) error {
	config := core.GetConfig()

	handler := Handler()

	addr := fmt.Sprintf(":%v", config.Port)
	logger.Info("Listening", zap.String("addr", addr))
	return handler.Start(addr)
}

func Handler() *echo.Echo {
	e := echo.New()
	e.GET("/link", link)
	return e
}

func link(c echo.Context) error {
	return c.Redirect(301, getRedirectUrl())
}

func getRedirectUrl() string {
	redirectUrlOnce.Do(func() {
		redirectUrl = fmt.Sprintf(
			"https://%v/oauth/authorize?response_type=code&client_id=%v&redirect_uri=%v",
			core.GetConfig().TWHost,
			url.PathEscape("endiangroup%2Ftransferwiser"),
			url.PathEscape(core.GetConfig().TWLoginRedirect))
	})
	return redirectUrl
}
