package web

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"

	"github.com/cloudflare/cfssl/revoke"
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

func (s *server) Run(env string, port, letsencryptPort int, caCert string) error {
	handler := s.MainHandler()
	tlsServer := handler.TLSServer
	tlsServer.TLSConfig = new(tls.Config)

	// Automatic let's encrypt
	handler.AutoTLSManager.Cache = autocert.DirCache(".cache")
	tlsServer.TLSConfig.GetCertificate = handler.AutoTLSManager.GetCertificate
	go http.ListenAndServe(fmt.Sprintf(":%v", letsencryptPort), handler.AutoTLSManager.HTTPHandler(nil))

	certpool := x509.NewCertPool()
	if !certpool.AppendCertsFromPEM([]byte(caCert)) {
		return errors.New("error loading ca certificate")
	}

	tlsServer.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	tlsServer.TLSConfig.ClientCAs = certpool
	tlsServer.TLSConfig.NextProtos = append(tlsServer.TLSConfig.NextProtos, "h2")
	tlsServer.Addr = fmt.Sprintf(":%v", port)

	s.logger.Info("Listening", zap.String("addr", tlsServer.Addr))
	return handler.StartServer(tlsServer)
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
		req := c.Request()
		if req.TLS == nil {
			s.logger.Error("missing client certificate")
			return c.String(401, "Unauthorized")
		}
		if len(req.TLS.PeerCertificates) == 0 {
			s.logger.Error("missing client certificate")
			return c.String(401, "Unauthorized")
		}

		cert := req.TLS.PeerCertificates[0]
		revoked, ok := revoke.VerifyCertificate(cert)
		if revoked || !ok {
			s.logger.Error("client certificate is revoked or incorrect", zap.Bool("revoked", revoked), zap.Bool("ok", ok))
			return c.String(401, "Unauthorized")
		}

		cn := cert.Subject.CommonName
		if cn == "" {
			s.logger.Error("client certificate doesn't contain CommonName (CN)", zap.String("subject", cert.Subject.String()))
			return c.String(401, "Unauthorized")
		}
		c.Set("username", cn)

		return next(c)
	}
}
