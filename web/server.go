package web

import (
	"fmt"

	"github.com/endiangroup/transferwiser/core"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type server struct {
	logger             *zap.Logger
	twAuth             core.TransferwiseAuthenticator
	twTransferProvider core.TransferwiseTransfersProvider
}

func NewServer(logger *zap.Logger, twAuth core.TransferwiseAuthenticator, twTransferProvider core.TransferwiseTransfersProvider) *server {
	return &server{
		logger:             logger,
		twAuth:             twAuth,
		twTransferProvider: twTransferProvider,
	}
}

func (s *server) Run(port int) error {
	handler := s.MainHandler()

	addr := fmt.Sprintf(":%v", port)
	s.logger.Info("Listening", zap.String("addr", addr))
	return handler.Start(addr)
}

func (s *server) MainHandler() *echo.Echo {
	e := echo.New()
	e.GET("/link", s.link)
	e.GET("/callback", s.callback)
	return e
}

func (s *server) link(c echo.Context) error {
	return c.Redirect(301, s.twAuth.RedirectUrl())
}

func (s *server) callback(c echo.Context) error {
	transferwiseRefreshToken := c.QueryParam("code")
	twData, err := s.twAuth.RefreshToken(transferwiseRefreshToken)
	if err != nil {
		s.logger.Error("error refreshing transferwise code", zap.Error(err))
		return echo.NewHTTPError(500, "error confirming transferwise authentication")
	}
	err = s.twTransferProvider.UseAuthentication(twData)
	if err != nil {
		s.logger.Error("error using autnehtication", zap.Error(err))
		return echo.NewHTTPError(500, "error configuring transferwise authentication")
	}
	return c.Redirect(301, "/")
}
