package transferwiser

import (
	"github.com/DATA-DOG/godog"
)

func theServiceHasNotBeenAuthenticatedWithTransferwise() error {
	return godog.ErrPending
}

func iVisitTheLinkToConnectWithTransferwise() error {
	return godog.ErrPending
}

func iShouldBeRedirectedToTheTransferwiseOauthLoginPage() error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^The service has not been authenticated with Transferwise$`, theServiceHasNotBeenAuthenticatedWithTransferwise)
	s.Step(`^I visit the link to connect with Transferwise$`, iVisitTheLinkToConnectWithTransferwise)
	s.Step(`^I should be redirected to the Transferwise Oauth login page$`, iShouldBeRedirectedToTheTransferwiseOauthLoginPage)
}
