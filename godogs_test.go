package transferwiser

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/endiangroup/transferwiser/core"
	"github.com/endiangroup/transferwiser/web"
	"github.com/labstack/echo"
)

type BDDContext struct {
	webApp *echo.Echo

	lastResponseRecorder *httptest.ResponseRecorder
}

func NewBDDContext() *BDDContext {
	return &BDDContext{
		webApp: web.Handler(),
	}
}

func (ctx *BDDContext) Init() {
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

	s.Step(`^The service has not been authenticated with Transferwise$`, ctx.TheServiceHasNotBeenAuthenticatedWithTransferwise)
	s.Step(`^I visit the link to connect with Transferwise$`, ctx.IVisitTheLinkToConnectWithTransferwise)
	s.Step(`^I should be redirected to the Transferwise Oauth login page$`, ctx.IShouldBeRedirectedToTheTransferwiseOauthLoginPage)
}

//
// Step definitions
//

func (ctx *BDDContext) TheServiceHasNotBeenAuthenticatedWithTransferwise() error {
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
	redirectUrl := ctx.lastResponseRecorder.HeaderMap.Get("location")
	partsToContain := []string{
		core.GetConfig().TWHost,
		"/oauth/authorize?",
		"response_type=code",
		fmt.Sprintf("client_id=%v", url.PathEscape("endiangroup%2Ftransferwiser")),
		fmt.Sprintf("redirect_uri=%v", url.PathEscape(core.GetConfig().TWLoginRedirect)),
	}
	if err := containsAll(redirectUrl, partsToContain...); err != nil {
		return err
	}
	return nil
}

func containsAll(s string, contained ...string) error {
	for _, check := range contained {
		if !strings.Contains(s, check) {
			return fmt.Errorf("'%v' expected to contain '%v'", s, check)
		}
	}
	return nil
}
