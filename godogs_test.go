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
	coreMocks "github.com/endiangroup/transferwiser/core/mocks"
	kvMocks "github.com/endiangroup/transferwiser/keyvalue/mocks"
	"github.com/endiangroup/transferwiser/web"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type BDDContext struct {
	transferwiseAPI     *coreMocks.TransferwiseAPI
	authStore           *kvMocks.Value
	transferwiseService *core.TransferwiseService
	webApp              *echo.Echo

	isAuthenticated bool

	lastResponseRecorder *httptest.ResponseRecorder
}

func NewBDDContext() *BDDContext {
	return &BDDContext{}
}

var redirectUrl = "https://sandbox.transferwise.tech/oauth/authorize?response_type=code&client_id=clientID&redirect_uri=http:%2F%2Fredirect%2Fhere"

func (ctx *BDDContext) Init() {
	ctx.transferwiseAPI = &coreMocks.TransferwiseAPI{}
	ctx.authStore = &kvMocks.Value{}
	ctx.transferwiseService = core.NewTransferwiseService(
		ctx.transferwiseAPI,
		"clientID",
		"sandbox.transferwise.tech",
		"http://redirect/here",
		ctx.authStore)

	webServer := web.NewServer(zap.NewNop(), ctx.transferwiseService)

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
	if ctx.transferwiseService.IsAuthenticated() {
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
	if receivedRedirectURL != redirectUrl {
		return fmt.Errorf("Redirect URL is %q, but returned %q", receivedRedirectURL, redirectUrl)
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
	ctx.transferwiseAPI.On("RefreshToken", token).Return(refreshTokenData, nil)
	ctx.authStore.On("PutString", refreshTokenData.AccessToken).Return(nil)
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
	if !ctx.transferwiseService.IsAuthenticated() {
		return errors.New("service expected to be authenticated, but it's not")
	}
	return nil
}
