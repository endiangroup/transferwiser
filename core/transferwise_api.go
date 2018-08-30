package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	DatetimeFormat = "2006-01-02 15:04:05"
)

//go:generate mockery -name=TransferwiseAPI
type TransferwiseAPI interface {
	Transfers() ([]*Transfer, error)
}

type transferwiseAPI struct {
	host     string
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

type transferwiseTransfer struct {
	ID             int64   `json:"id"`
	Created        string  `json:"created"`
	Status         string  `json:"status"`
	SourceValue    float64 `json:"sourceValue"`
	SourceCurrency string  `json:"sourceCurrency"`
	TargetValue    float64 `json:"targetValue"`
	TargetCurrency string  `json:"targetCurrency"`
	Rate           float64 `json:"rate"`
}

func (twTransfer *transferwiseTransfer) toTransfer() (*Transfer, error) {
	created, err := time.Parse(DatetimeFormat, twTransfer.Created)
	if err != nil {
		return nil, err
	}
	transfer := &Transfer{
		ID:             twTransfer.ID,
		CreatedAt:      created,
		Status:         twTransfer.Status,
		SourceValue:    twTransfer.SourceValue,
		SourceCurrency: twTransfer.SourceCurrency,
		TargetValue:    twTransfer.TargetValue,
		TargetCurrency: twTransfer.TargetCurrency,
		ExchangeRate:   twTransfer.Rate,
	}

	return transfer, nil
}

func NewTransferwiseAPI(host, token string) *transferwiseAPI {
	return &transferwiseAPI{
		host:     host,
		apiToken: token,
	}
}

func (tw *transferwiseAPI) Transfers() ([]*Transfer, error) {
	url := fmt.Sprintf("https://%v/v1/transfers?offset=0&limit=20", tw.host)
	httpClient := http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error preparing transferwise request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", tw.apiToken))
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error sending transferwise request")
	}

	defer res.Body.Close()
	twTransfers := []*transferwiseTransfer{}
	err = json.NewDecoder(res.Body).Decode(&twTransfers)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding transferwise request")
	}

	transfers := make([]*Transfer, len(twTransfers))
	for i, twTransfer := range twTransfers {
		transfers[i], err = twTransfer.toTransfer()
		if err != nil {
			return nil, errors.Wrap(err, "error parsing transferwise request")
		}
	}

	return transfers, nil
}
