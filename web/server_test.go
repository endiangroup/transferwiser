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

var redirectUrl = "https://sandbox.transferwise.tech/oauth/authorize?response_type=code&client_id=clientID&redirect_uri=http:%2F%2Fredirect%2Fhere"

func TestTransferwiseLinkRedirect(t *testing.T) {
	transferwiseAPI := &coreMocks.TransferwiseAPI{}
	authStore := &kvMocks.Value{}
	transferwiseService := core.NewTransferwiseService(
		transferwiseAPI,
		"clientID",
		"sandbox.transferwise.tech",
		"http://redirect/here",
		authStore)

	webServer := NewServer(zap.NewNop(), transferwiseService)

	webApp := webServer.MainHandler()

	req, err := http.NewRequest(echo.GET, "/oauth/link", nil)
	require.NoError(t, err)

	resRec := httptest.NewRecorder()
	webApp.ServeHTTP(resRec, req)

	require.Equal(t, 301, resRec.Code)
	require.Equal(t, redirectUrl, resRec.HeaderMap.Get("location"))
}

func TestTransferwiseCallback(t *testing.T) {
	token := "asdf1234"

	transferwiseAPI := &coreMocks.TransferwiseAPI{}
	authStore := &kvMocks.Value{}
	transferwiseService := core.NewTransferwiseService(
		transferwiseAPI,
		"clientID",
		"sandbox.transferwise.tech",
		"http://redirect/here",
		authStore)

	webServer := NewServer(zap.NewNop(), transferwiseService)

	webApp := webServer.MainHandler()

	refreshTokenData := &core.RefreshTokenData{
		AccessToken:  "myaccesstoken",
		TokenType:    "bearer",
		RefreshToken: "abcd1234",
		ExpiresIn:    3600,
		Scope:        "transfers",
	}
	transferwiseAPI.On("RefreshToken", token).Return(refreshTokenData, nil)
	authStore.On("PutString", refreshTokenData.AccessToken).Return(nil)

	req, err := http.NewRequest(echo.GET, fmt.Sprintf("/oauth/callback?code=%v", token), nil)
	require.NoError(t, err)

	resRec := httptest.NewRecorder()

	require.False(t, transferwiseService.IsAuthenticated())
	webApp.ServeHTTP(resRec, req)
	require.True(t, transferwiseService.IsAuthenticated())
}
