// +build integration

package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/endiangroup/transferwiser/core"
	coreMocks "github.com/endiangroup/transferwiser/core/mocks"
	kvMocks "github.com/endiangroup/transferwiser/keyvalue/mocks"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testData struct {
	transferwiseAPI     *coreMocks.TransferwiseAPI
	authStore           *kvMocks.Value
	transferwiseService *core.TransferwiseService
	webApp              *echo.Echo
}

func getTestData(t *testing.T) *testData {
	transferwiseAPI := &coreMocks.TransferwiseAPI{}
	authStore := &kvMocks.Value{}
	transferwiseService := core.NewTransferwiseService(
		transferwiseAPI,
		"clientID",
		"sandbox.transferwise.tech",
		"http://redirect/here",
		authStore)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	webServer := NewServer(logger, transferwiseService)

	webApp := webServer.MainHandler()
	return &testData{
		transferwiseAPI:     transferwiseAPI,
		authStore:           authStore,
		transferwiseService: transferwiseService,
		webApp:              webApp,
	}
}

func TestTransferwiseLinkRedirect_Redirects(t *testing.T) {
	data := getTestData(t)

	req, err := http.NewRequest(echo.GET, "/oauth/link", nil)
	require.NoError(t, err)
	req.Header.Set("DN", "EMAIL=admin@endian.io,CN=admin,C=UK")
	req.Header.Set("VERIFIED", "SUCCESS")

	resRec := httptest.NewRecorder()
	data.webApp.ServeHTTP(resRec, req)

	require.Equal(t, 301, resRec.Code)
	require.Equal(t, data.transferwiseService.RedirectUrl(), resRec.HeaderMap.Get("location"))
}

func TestTransferwiseLinkRedirect_RequiresAdmin(t *testing.T) {
	data := getTestData(t)

	req, err := http.NewRequest(echo.GET, "/oauth/link", nil)
	require.NoError(t, err)
	req.Header.Set("DN", "EMAIL=accounting@endian.io,CN=accountant1,C=UK")
	req.Header.Set("VERIFIED", "SUCCESS")

	resRec := httptest.NewRecorder()
	data.webApp.ServeHTTP(resRec, req)

	require.Equal(t, 403, resRec.Code)
}

func TestTransferwiseLinkRedirect_RequiresVerifiedCertificate(t *testing.T) {
	data := getTestData(t)

	req, err := http.NewRequest(echo.GET, "/oauth/link", nil)
	require.NoError(t, err)
	req.Header.Set("DN", "EMAIL=admin@endian.io,CN=admin,C=UK")
	req.Header.Set("VERIFIED", "NONE")

	resRec := httptest.NewRecorder()
	data.webApp.ServeHTTP(resRec, req)

	require.Equal(t, 401, resRec.Code)
}

func TestTransferwiseCallback_UsesCredentials(t *testing.T) {
	token := "asdf1234"
	data := getTestData(t)

	refreshTokenData := &core.RefreshTokenData{
		AccessToken:  "myaccesstoken",
		TokenType:    "bearer",
		RefreshToken: "abcd1234",
		ExpiresIn:    3600,
		Scope:        "transfers",
	}
	data.transferwiseAPI.On("RefreshToken", token).Return(refreshTokenData, nil)
	data.authStore.On("PutString", refreshTokenData.AccessToken).Return(nil)

	req, err := http.NewRequest(echo.GET, fmt.Sprintf("/oauth/callback?code=%v", token), nil)
	require.NoError(t, err)
	req.Header.Set("DN", "EMAIL=admin@endian.io,CN=admin,C=UK")
	req.Header.Set("VERIFIED", "SUCCESS")

	resRec := httptest.NewRecorder()

	require.False(t, data.transferwiseService.IsAuthenticated())
	data.webApp.ServeHTTP(resRec, req)
	require.Equal(t, 301, resRec.Code)
	require.Equal(t, "/", resRec.HeaderMap.Get("location"))
	require.True(t, data.transferwiseService.IsAuthenticated())
}

func TestTransferwiseCallback_RequiresAdmin(t *testing.T) {
	token := "asdf1234"
	data := getTestData(t)

	req, err := http.NewRequest(echo.GET, fmt.Sprintf("/oauth/callback?code=%v", token), nil)
	require.NoError(t, err)
	req.Header.Set("DN", "EMAIL=accounting@endian.io,CN=accountant1,C=UK")
	req.Header.Set("VERIFIED", "SUCCESS")

	resRec := httptest.NewRecorder()

	data.webApp.ServeHTTP(resRec, req)
	require.Equal(t, 403, resRec.Code)
	require.False(t, data.transferwiseService.IsAuthenticated())
}

func TestTransferwiseCallback_RequiresVerifiedCertificate(t *testing.T) {
	token := "asdf1234"
	data := getTestData(t)

	req, err := http.NewRequest(echo.GET, fmt.Sprintf("/oauth/callback?code=%v", token), nil)
	require.NoError(t, err)
	req.Header.Set("DN", "EMAIL=admin@endian.io,CN=admin,C=UK")
	req.Header.Set("VERIFIED", "NONE")

	resRec := httptest.NewRecorder()

	data.webApp.ServeHTTP(resRec, req)
	require.Equal(t, 401, resRec.Code)
	require.False(t, data.transferwiseService.IsAuthenticated())
}
