package core

//go:generate mockery -name=TransferwiseAPI
type TransferwiseAPI interface {
	Transfers() ([]*Transfer, error)
	RefreshToken(code string) (*RefreshTokenData, error)
}

type RefreshTokenData struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
}

type transferwiseAPI struct {
}

type Transfer struct {
}

func NewTransferwiseAPI() *transferwiseAPI {
	return &transferwiseAPI{}
}

func (tw *transferwiseAPI) Transfers() ([]*Transfer, error) {
	return []*Transfer{}, nil
}

func (tw *transferwiseAPI) RefreshToken(code string) (*RefreshTokenData, error) {
	return &RefreshTokenData{}, nil
}
