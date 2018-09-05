package web

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/endiangroup/transferwiser/core"
	"github.com/gocarina/gocsv"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
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
	tlsServer := handler.TLSServer
	tlsServer.TLSConfig = new(tls.Config)

	// Automatic let's encrypt
	handler.AutoTLSManager.Cache = autocert.DirCache(".cache")
	tlsServer.TLSConfig.GetCertificate = handler.AutoTLSManager.GetCertificate
	go http.ListenAndServe(":http", handler.AutoTLSManager.HTTPHandler(nil))

	certpool := x509.NewCertPool()
	if !certpool.AppendCertsFromPEM(cert) {
		return errors.New("error loading ca certificate")
	}

	tlsServer.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	tlsServer.TLSConfig.ClientCAs = certpool
	tlsServer.TLSConfig.NextProtos = append(tlsServer.TLSConfig.NextProtos, "h2")
	tlsServer.Addr = fmt.Sprintf(":%v", port)

	s.logger.Info("Listening", zap.String("addr", tlsServer.Addr))
	return handler.StartServer(handler.TLSServer)
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
