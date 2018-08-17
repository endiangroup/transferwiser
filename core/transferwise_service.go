package core

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/endiangroup/transferwiser/keyvalue"
)

type TransferwiseService struct {
	Api                   TransferwiseAPI
	clientId              string
	transferwiseHost      string
	loginRedirectCallback string
	authorizationStore    keyvalue.Value

	formatRedirectUrlOnce sync.Once
	redirectUrl           string

	authLock        sync.Mutex
	isAuthenticated bool
}

func NewTransferwiseService(api TransferwiseAPI, clientId, transferwiseHost, loginRedirectCallback string, authorizationStore keyvalue.Value) *TransferwiseService {
	return &TransferwiseService{
		Api:                   api,
		clientId:              clientId,
		transferwiseHost:      transferwiseHost,
		loginRedirectCallback: loginRedirectCallback,
		authorizationStore:    authorizationStore,
		isAuthenticated:       false,
	}
}

func (tw *TransferwiseService) RedirectUrl() string {
	tw.formatRedirectUrlOnce.Do(func() {
		tw.redirectUrl = fmt.Sprintf(
			"https://%v/oauth/authorize?response_type=code&client_id=%v&redirect_uri=%v",
			tw.transferwiseHost,
			url.PathEscape(tw.clientId),
			url.PathEscape(tw.loginRedirectCallback))
	})
	return tw.redirectUrl
}

func (tw *TransferwiseService) UseAuthentication(data *RefreshTokenData) error {
	tw.authLock.Lock()
	defer tw.authLock.Unlock()
	if err := tw.authorizationStore.PutString(data.AccessToken); err != nil {
		return err
	}
	tw.isAuthenticated = true
	return nil
}

func (tw *TransferwiseService) IsAuthenticated() bool {
	tw.authLock.Lock()
	defer tw.authLock.Unlock()
	return tw.isAuthenticated
}
