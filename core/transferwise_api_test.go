package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	gock "gopkg.in/h2non/gock.v1"
)

func TestTransfers_SimpleList(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution

	data := `[
		{
			"id": 471041,
			"user": null,
			"targetAccount": 14177602,
			"sourceAccount": null,
			"quote": null,
			"status": "incoming_payment_waiting",
			"reference": "",
			"rate": 1.1119,
			"created": "2018-08-29 17:58:24",
			"business": 724,
			"transferRequest": null,
			"details": {
				"reference": ""
			},
			"hasActiveIssues": false,
			"sourceCurrency": "GBP",
			"sourceValue": 29894.57,
			"targetCurrency": "EUR",
			"targetValue": 33239.77
		}
	]`

	gock.New(fmt.Sprintf("https://%v", GetConfig().TwHost)).
		Get("/v1/transfers").
		MatchParam("offset", "0").
		MatchParam("limit", "20").
		MatchHeader("Authorization", fmt.Sprintf("Bearer %v", GetConfig().TwAPIToken)).
		Reply(200).
		Type("application/json").
		BodyString(data)

	api := NewTransferwiseAPI(GetConfig().TwHost, GetConfig().TwAPIToken)
	transfers, err := api.Transfers()
	require.NoError(t, err)
	require.Len(t, transfers, 1)

	transfer := transfers[0]
	require.Equal(t, int64(471041), transfer.ID)
	require.Equal(t, time.Date(2018, 8, 29, 17, 58, 24, 0, time.UTC), transfer.CreatedAt)
	require.Equal(t, "incoming_payment_waiting", transfer.Status)
	// require.Equal(t, , transfer.RecipientName)
	require.Equal(t, 29894.57, transfer.SourceValue)
	require.Equal(t, "GBP", transfer.SourceCurrency)
	require.Equal(t, 33239.77, transfer.TargetValue)
	require.Equal(t, "EUR", transfer.TargetCurrency)
	// require.Equal(t, , transfer.Fee)
	require.Equal(t, 1.1119, transfer.ExchangeRate)
}
