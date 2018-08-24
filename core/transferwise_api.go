package core

//go:generate mockery -name=TransferwiseAPI
type TransferwiseAPI interface {
	Transfers() ([]*Transfer, error)
}

type transferwiseAPI struct {
	apiToken string
}

type Transfer struct {
}

func NewTransferwiseAPI(token string) *transferwiseAPI {
	return &transferwiseAPI{
		apiToken: token,
	}
}

func (tw *transferwiseAPI) Transfers() ([]*Transfer, error) {
	return []*Transfer{}, nil
}
