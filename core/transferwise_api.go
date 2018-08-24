package core

import "time"

//go:generate mockery -name=TransferwiseAPI
type TransferwiseAPI interface {
	Transfers() ([]*Transfer, error)
}

type transferwiseAPI struct {
	apiToken string
}

type Transfer struct {
	ID             int64     `csv:"id"`
	CreatedAt      time.Time `csv:"created_at"`
	Status         string    `csv:"status"`
	RecipientName  string    `csv:"recipient_name"`
	SourceValue    float64   `csv:"source_value"`
	SourceCurrency string    `csv:"source_currency"`
	TargetValue    float64   `csv:"target_value"`
	TargetCurrency string    `csv:"target_currency"`
	Fee            float64   `csv:"fee"`
	ExchangeRate   float64   `csv:"exchange_rate"`
}

func NewTransferwiseAPI(token string) *transferwiseAPI {
	return &transferwiseAPI{
		apiToken: token,
	}
}

func (tw *transferwiseAPI) Transfers() ([]*Transfer, error) {
	return []*Transfer{}, nil
}
