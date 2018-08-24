// +build integration

package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	coreMocks "github.com/endiangroup/transferwiser/core/mocks"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testData struct {
	transferwiseAPI *coreMocks.TransferwiseAPI
	webApp          *echo.Echo
}

func getTestData(t *testing.T) *testData {
	transferwiseAPI := &coreMocks.TransferwiseAPI{}

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	webServer := NewServer(logger, transferwiseAPI)

	webApp := webServer.MainHandler()
	return &testData{
		transferwiseAPI: transferwiseAPI,
		webApp:          webApp,
	}
}

func TestTransferwiseTransfers_RespondsWithACsv(t *testing.T) {
	data := getTestData(t)

	req, err := http.NewRequest(echo.GET, "/transfers.csv", nil)
	require.NoError(t, err)
	req.Header.Set("DN", "EMAIL=admin@endian.io,CN=admin,C=UK")
	req.Header.Set("VERIFIED", "SUCCESS")

	resRec := httptest.NewRecorder()
	data.webApp.ServeHTTP(resRec, req)

	require.Equal(t, 200, resRec.Code)
	require.Equal(t, "text/csv", resRec.Header().Get(echo.HeaderContentType))
}

func TestTransferwiseTransfers_RequiresCN(t *testing.T) {
	data := getTestData(t)

	req, err := http.NewRequest(echo.GET, "/transfers.csv", nil)
	require.NoError(t, err)
	req.Header.Set("DN", "EMAIL=admin@endian.io,C=UK")
	req.Header.Set("VERIFIED", "SUCCESS")

	resRec := httptest.NewRecorder()
	data.webApp.ServeHTTP(resRec, req)

	require.Equal(t, 401, resRec.Code)
}

func TestTransferwiseTransfers_RequiresVerifiedCertificate(t *testing.T) {
	data := getTestData(t)

	req, err := http.NewRequest(echo.GET, "/transfers.csv", nil)
	require.NoError(t, err)
	req.Header.Set("DN", "EMAIL=admin@endian.io,CN=admin,C=UK")
	req.Header.Set("VERIFIED", "NONE")

	resRec := httptest.NewRecorder()
	data.webApp.ServeHTTP(resRec, req)

	require.Equal(t, 401, resRec.Code)
}
