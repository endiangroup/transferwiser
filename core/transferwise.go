package core

import (
	"fmt"
	"net/url"
	"sync"
)

//go:generate mockery -name=TransferwiseAuthenticator
type TransferwiseAuthenticator interface {
	RedirectUrl() string
	RefreshToken(code string) (*RefreshTokenData, error)
}

//go:generate mockery -name=TransferwiseTransfersProvider
type TransferwiseTransfersProvider interface {
	Transfers() ([]*Transfer, error)
	UseAuthentication(*RefreshTokenData) error
	IsAuthenticated() bool
}

type RefreshTokenData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

type Transfer struct {
}

type TransferwiseService struct {
	clientId              string
	transferwiseHost      string
	loginRedirectCallback string

	formatRedirectUrlOnce sync.Once
	redirectUrl           string
}

func NewTransferwiseService(clientId, transferwiseHost, loginRedirectCallback string) *TransferwiseService {
	return &TransferwiseService{
		clientId:              clientId,
		transferwiseHost:      transferwiseHost,
		loginRedirectCallback: loginRedirectCallback,
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

func (tw *TransferwiseService) AccessTokenFromcCode(code string) (string, error) {
	return "", nil
}

func (tw *TransferwiseService) UseAuthentication(*RefreshTokenData) error {
	return nil
}

func (tw *TransferwiseService) IsAuthenticated() bool {
	return false
}

func (tw *TransferwiseService) Transfers() ([]*Transfer, error) {
	return []*Transfer{}, nil
}
