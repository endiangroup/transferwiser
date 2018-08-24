package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/endiangroup/transferwiser/core"
	"github.com/gocarina/gocsv"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type server struct {
	logger          *zap.Logger
	transferwiseAPI core.TransferwiseAPI
}

func NewServer(logger *zap.Logger, transferwiseAPI core.TransferwiseAPI) *server {
	return &server{
		logger:          logger,
		transferwiseAPI: transferwiseAPI,
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
	e.GET("/transfers.csv", s.transfersCSV)
	return e
}

func (s *server) transfersCSV(c echo.Context) error {
	transfers, err := s.transferwiseAPI.Transfers()
	if err != nil {
		s.logger.Error("error fetching transferwise transfers", zap.Error(err))
		return c.String(500, "Internal server error")
	}
	c.Response().Header().Set(echo.HeaderContentType, "text/csv")
	c.Response().WriteHeader(http.StatusOK)

	return gocsv.Marshal(&transfers, c.Response())
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
		if cn == "" {
			s.logger.Error("DN header don't contain DN", zap.String("dn", dnHeader))
			return c.String(401, "Unauthorized")
		}
		c.Set("username", cn)

		return next(c)
	}
}
