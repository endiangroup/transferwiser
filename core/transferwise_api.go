package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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
	RecipientID    int64     `json:"-"`
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
	TargetAccount  int64   `json:"targetAccount"`
	Rate           float64 `json:"rate"`
}

type transferwiseAccount struct {
	Name string `json:"accountHolderName"`
}

func (twTransfer *transferwiseTransfer) toTransfer() (*Transfer, error) {
	created, err := time.Parse(DatetimeFormat, twTransfer.Created)
	if err != nil {
		return nil, err
	}
	transfer := &Transfer{
		ID:             twTransfer.ID,
		CreatedAt:      created,
		RecipientID:    twTransfer.TargetAccount,
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
	createdDateStart := getFirstOfTwoMonthsAgo(time.Now())
	httpClient := http.Client{
		Timeout: time.Second * 20,
	}
	twTransfers, err := tw.getTransfers(httpClient, 0, 20, createdDateStart)
	if err != nil {
		return nil, errors.Wrap(err, "error getting transferwise transfers")
	}

	accounts := map[int64]bool{}
	transfers := make([]*Transfer, len(twTransfers))
	for i, twTransfer := range twTransfers {
		transfers[i], err = twTransfer.toTransfer()
		accounts[twTransfer.TargetAccount] = true
		if err != nil {
			return nil, errors.Wrap(err, "error parsing transferwise transfers request")
		}
	}
	names, err := tw.getAccountNames(httpClient, accounts)
	if err != nil {
		return nil, errors.Wrap(err, "error getting account names")
	}

	for _, transfer := range transfers {
		transfer.RecipientName = names[transfer.RecipientID]
	}

	return transfers, nil
}

func (tw *transferwiseAPI) getTransfers(httpClient http.Client, offset, limit int, createdDateStart string) ([]*transferwiseTransfer, error) {
	url := fmt.Sprintf("https://%v/v1/transfers?offset=%d&limit=%d&createdDateStart=%v", tw.host, offset, limit, createdDateStart)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error preparing transferwise transfers request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", tw.apiToken))
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error sending transferwise transfers request")
	}

	defer res.Body.Close()
	twTransfers := []*transferwiseTransfer{}
	err = json.NewDecoder(res.Body).Decode(&twTransfers)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding transferwise transfers request")
	}
	return twTransfers, nil
}

func (tw *transferwiseAPI) getAccountNames(httpClient http.Client, idsSet map[int64]bool) (map[int64]string, error) {
	names := map[int64]string{}
	var g errgroup.Group
	var mtx sync.Mutex

	for id, _ := range idsSet {
		_id := id
		url := fmt.Sprintf("https://%v/v1/accounts/%v", tw.host, id)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.Wrap(err, "error preparing transferwise account request")
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", tw.apiToken))

		g.Go(func() error {
			res, err := httpClient.Do(req)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			twAccount := &transferwiseAccount{}
			err = json.NewDecoder(res.Body).Decode(&twAccount)
			if err != nil {
				return errors.Wrap(err, "error decoding transferwise account request")
			}
			mtx.Lock()
			names[_id] = twAccount.Name
			mtx.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return names, nil
}

func getFirstOfTwoMonthsAgo(t time.Time) string {
	twoMonthsAgo := t.AddDate(0, -2, 0)
	return fmt.Sprintf("%v-%v-%v", twoMonthsAgo.Year(), int(twoMonthsAgo.Month()), 1)
}
