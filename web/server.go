package web

import (
	"fmt"

	"github.com/endiangroup/transferwiser/core"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type server struct {
	logger              *zap.Logger
	transferwiseService *core.TransferwiseService
}

func NewServer(logger *zap.Logger, transferwiseService *core.TransferwiseService) *server {
	return &server{
		logger:              logger,
		transferwiseService: transferwiseService,
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
	return c.Redirect(301, s.transferwiseService.RedirectUrl())
}

func (s *server) callback(c echo.Context) error {
	transferwiseRefreshToken := c.QueryParam("code")
	twData, err := s.transferwiseService.Api.RefreshToken(transferwiseRefreshToken)
	if err != nil {
		s.logger.Error("error refreshing transferwise code", zap.Error(err))
		return echo.NewHTTPError(500, "error confirming transferwise authentication")
	}
	err = s.transferwiseService.UseAuthentication(twData)
	if err != nil {
		s.logger.Error("error using authentication", zap.Error(err))
		return echo.NewHTTPError(500, "error configuring transferwise authentication")
	}
	return c.Redirect(301, "/")
}
