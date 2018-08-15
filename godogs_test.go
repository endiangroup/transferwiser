package transferwiser

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/endiangroup/transferwiser/core"
	"github.com/endiangroup/transferwiser/core/mocks"
	"github.com/endiangroup/transferwiser/web"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

var (
	redirectURL = "http://must-redirect-here.com"
)

type BDDContext struct {
	webAuthenticator *mocks.TransferwiseAuthenticator
	transferProvider *mocks.TransferwiseTransfersProvider
	webApp           *echo.Echo

	isAuthenticated bool

	lastResponseRecorder *httptest.ResponseRecorder
}

func NewBDDContext() *BDDContext {
	return &BDDContext{}
}

func (ctx *BDDContext) Init() {
	authMock := &mocks.TransferwiseAuthenticator{}
	transferProviderMock := &mocks.TransferwiseTransfersProvider{}
	ctx.isAuthenticated = false
	webServer := web.NewServer(zap.NewNop(), authMock, transferProviderMock)

	authMock.On("RedirectUrl").Return(redirectURL)
	transferProviderMock.On("IsAuthenticated").Return(func() bool {
		return ctx.isAuthenticated
	})
	transferProviderMock.On("UseAuthentication", mock.Anything).Return(func(data *core.RefreshTokenData) error {
		ctx.isAuthenticated = true
		return nil
	})

	ctx.webAuthenticator = authMock
	ctx.transferProvider = transferProviderMock
	ctx.webApp = webServer.MainHandler()
}

func (ctx *BDDContext) Teardown() {
}

func FeatureContext(s *godog.Suite) {
	ctx := NewBDDContext()

	s.BeforeScenario(func(interface{}) {
		ctx.Init()
	})
	s.AfterScenario(func(interface{}, error) {
		ctx.Teardown()
	})

	s.Step(`^the service has not been authenticated with Transferwise$`, ctx.TheServiceHasNotBeenAuthenticatedWithTransferwise)
	s.Step(`^I visit the link to connect with Transferwise$`, ctx.IVisitTheLinkToConnectWithTransferwise)
	s.Step(`^I should be redirected to the Transferwise authorization login page$`, ctx.IShouldBeRedirectedToTheTransferwiseOauthLoginPage)
	s.Step(`^Transferwise refresh token response for \'(\w+)\' is:$`, ctx.TransferwiseRefreshTokenResponseIs)
	s.Step(`^I return from Transferwise OAuth with code \'(\w+)\'$`, ctx.IReturnFromTransferwiseOAuthWithCode)
	s.Step(`^the service is authenticated with Transferwise$`, ctx.TheServiceIsAuthenticatedWithTransferwise)
}

//
// Step definitions
//

func (ctx *BDDContext) TheServiceHasNotBeenAuthenticatedWithTransferwise() error {
	if ctx.transferProvider.IsAuthenticated() {
		return errors.New("service expected to be authenticated, but it's not")
	}
	return nil
}

func (ctx *BDDContext) IVisitTheLinkToConnectWithTransferwise() error {
	req, err := http.NewRequest("GET", "/link", nil)
	if err != nil {
		return err
	}
	resRec := httptest.NewRecorder()
	ctx.webApp.ServeHTTP(resRec, req)

	ctx.lastResponseRecorder = resRec
	return nil
}

func (ctx *BDDContext) IShouldBeRedirectedToTheTransferwiseOauthLoginPage() error {
	if ctx.lastResponseRecorder.Code != 301 {
		return fmt.Errorf("expected 301 redirect, got %v", ctx.lastResponseRecorder.Code)
	}
	receivedRedirectURL := ctx.lastResponseRecorder.HeaderMap.Get("location")
	if receivedRedirectURL != redirectURL {
		return fmt.Errorf("Redirect URL is %q, but returned %q", receivedRedirectURL, redirectURL)
	}
	return nil
}

func (ctx *BDDContext) TransferwiseRefreshTokenResponseIs(token string, data *gherkin.DataTable) error {
	refreshTokenData := &core.RefreshTokenData{}
	for _, row := range data.Rows {
		switch row.Cells[0].Value {
		case "access_token":
			refreshTokenData.AccessToken = row.Cells[1].Value
		case "token_type":
			refreshTokenData.TokenType = row.Cells[1].Value
		case "refresh_token":
			refreshTokenData.RefreshToken = row.Cells[1].Value
		case "expires_in":
			num, err := strconv.Atoi(row.Cells[1].Value)
			if err != nil {
				return err
			}
			refreshTokenData.ExpiresIn = int64(num)
		case "scope":
			refreshTokenData.Scope = row.Cells[1].Value

		}
	}
	ctx.webAuthenticator.On("RefreshToken", token).Return(refreshTokenData, nil)
	return nil
}

func (ctx *BDDContext) IReturnFromTransferwiseOAuthWithCode(token string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("/callback?code=%v", token), nil)
	if err != nil {
		return err
	}
	resRec := httptest.NewRecorder()
	ctx.webApp.ServeHTTP(resRec, req)

	return nil
}

func (ctx *BDDContext) TheServiceIsAuthenticatedWithTransferwise() error {
	if !ctx.transferProvider.IsAuthenticated() {
		return errors.New("service expected to be authenticated, but it's not")
	}
	return nil
}
