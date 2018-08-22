package web

import (
	"fmt"
	"strings"

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
	e.Use(s.authenticate)
	e.GET("/oauth/link", s.link, s.requireUsername("admin"))
	e.GET("/oauth/callback", s.callback, s.requireUsername("admin"))
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

func (s *server) authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		verified := c.Request().Header.Get("VERIFIED")
		if verified != "SUCCESS" {
			s.logger.Error("client certificate not verified", zap.String("verified", verified))
			return c.String(401, "Unauthorized")
		}
		dnHeader := c.Request().Header.Get("DN")
		if dnHeader == "" {
			s.logger.Error("headers don't contain DN")
			return c.String(401, "Unauthorized")
		}
		parts := strings.Split(dnHeader, ",")
		var cn string
		for _, part := range parts {
			if strings.HasPrefix(part, "CN=") {
				cn = strings.TrimPrefix(part, "CN=")
				break
			}
		}
		if dnHeader == "" {
			s.logger.Error("DN header don't contain DN", zap.String("dn", dnHeader))
			return c.String(401, "Unauthorized")
		}
		c.Set("username", cn)

		return next(c)
	}
}
