package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	gock "gopkg.in/h2non/gock.v1"
)

var mockedAccounts = map[int]string{
	1: `{
		"id": 1,
		"business": null,
		"profile": 723,
		"accountHolderName": "Neo",
		"currency": "EUR",
		"country": "DE",
		"type": "iban",
		"details": {
			"address": null,
			"email": "neo@matrix.com",
			"legalType": "PRIVATE",
			"accountNumber": null,
			"sortCode": null,
			"abartn": null,
			"accountType": null,
			"bankgiroNumber": null,
			"ifscCode": null,
			"bsbCode": null,
			"institutionNumber": null,
			"transitNumber": null,
			"phoneNumber": null,
			"bankCode": null,
			"russiaRegion": null,
			"routingNumber": null,
			"branchCode": null,
			"cpf": null,
			"cardNumber": null,
			"idType": null,
			"idNumber": null,
			"idCountryIso3": null,
			"idValidFrom": null,
			"idValidTo": null,
			"clabe": null,
			"swiftCode": null,
			"dateOfBirth": null,
			"clearingNumber": null,
			"bankName": null,
			"branchName": null,
			"businessNumber": null,
			"province": null,
			"city": null,
			"rut": null,
			"token": null,
			"cnpj": null,
			"payinReference": null,
			"pspReference": null,
			"orderId": null,
			"idDocumentType": null,
			"idDocumentNumber": null,
			"targetProfile": null,
			"iban": "DE89370400440532013000",
			"bic": null,
			"IBAN": "DE89370400440532013000",
			"BIC": null
		},
		"user": 12345656
	}`,
	2: `{
		"id": 2,
		"business": null,
		"profile": 723,
		"accountHolderName": "Morpheus",
		"currency": "EUR",
		"country": "DE",
		"type": "iban",
		"details": {
			"address": null,
			"email": "morpheus@matrix.com",
			"legalType": "PRIVATE",
			"accountNumber": null,
			"sortCode": null,
			"abartn": null,
			"accountType": null,
			"bankgiroNumber": null,
			"ifscCode": null,
			"bsbCode": null,
			"institutionNumber": null,
			"transitNumber": null,
			"phoneNumber": null,
			"bankCode": null,
			"russiaRegion": null,
			"routingNumber": null,
			"branchCode": null,
			"cpf": null,
			"cardNumber": null,
			"idType": null,
			"idNumber": null,
			"idCountryIso3": null,
			"idValidFrom": null,
			"idValidTo": null,
			"clabe": null,
			"swiftCode": null,
			"dateOfBirth": null,
			"clearingNumber": null,
			"bankName": null,
			"branchName": null,
			"businessNumber": null,
			"province": null,
			"city": null,
			"rut": null,
			"token": null,
			"cnpj": null,
			"payinReference": null,
			"pspReference": null,
			"orderId": null,
			"idDocumentType": null,
			"idDocumentNumber": null,
			"targetProfile": null,
			"iban": "DE89370400440532013000",
			"bic": null,
			"IBAN": "DE89370400440532013000",
			"BIC": null
		},
		"user": 12345656
	}`,
}

func TestTransfers_SimpleList(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.DisableNetworking()

	transfersData := `[
		{
			"id": 877677,
			"user": null,
			"targetAccount": 1,
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
		BodyString(transfersData)

	gock.New(fmt.Sprintf("https://%v", GetConfig().TwHost)).
		Get(fmt.Sprintf("/v1/accounts/%v", 1)).
		MatchHeader("Authorization", fmt.Sprintf("Bearer %v", GetConfig().TwAPIToken)).
		Reply(200).
		Type("application/json").
		BodyString(mockedAccounts[1])

	api := NewTransferwiseAPI(GetConfig().TwHost, GetConfig().TwAPIToken)
	transfers, err := api.Transfers()
	require.NoError(t, err)
	require.Len(t, transfers, 1)

	transfer := transfers[0]
	require.Equal(t, int64(877677), transfer.ID)
	require.Equal(t, time.Date(2018, 8, 29, 17, 58, 24, 0, time.UTC), transfer.CreatedAt)
	require.Equal(t, "incoming_payment_waiting", transfer.Status)
	require.Equal(t, "Neo", transfer.RecipientName)
	require.Equal(t, 29894.57, transfer.SourceValue)
	require.Equal(t, "GBP", transfer.SourceCurrency)
	require.Equal(t, 33239.77, transfer.TargetValue)
	require.Equal(t, "EUR", transfer.TargetCurrency)
	// require.Equal(t, , transfer.Fee)
	require.Equal(t, 1.1119, transfer.ExchangeRate)
}

func TestTransfers_MultipleTransfers(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	gock.DisableNetworking()

	transfersData := `[
		{
			"id": 877677,
			"user": null,
			"targetAccount": 1,
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
		},
		{
			"id": 877677,
			"user": null,
			"targetAccount": 2,
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
			"sourceValue": 2989.45,
			"targetCurrency": "EUR",
			"targetValue": 3323.97
		}
	]`

	gock.New(fmt.Sprintf("https://%v", GetConfig().TwHost)).
		Get("/v1/transfers").
		MatchParam("offset", "0").
		MatchParam("limit", "20").
		MatchHeader("Authorization", fmt.Sprintf("Bearer %v", GetConfig().TwAPIToken)).
		Reply(200).
		Type("application/json").
		BodyString(transfersData)

	gock.New(fmt.Sprintf("https://%v", GetConfig().TwHost)).
		Get(fmt.Sprintf("/v1/accounts/%v", 1)).
		MatchHeader("Authorization", fmt.Sprintf("Bearer %v", GetConfig().TwAPIToken)).
		Reply(200).
		Type("application/json").
		BodyString(mockedAccounts[1])

	gock.New(fmt.Sprintf("https://%v", GetConfig().TwHost)).
		Get(fmt.Sprintf("/v1/accounts/%v", 2)).
		MatchHeader("Authorization", fmt.Sprintf("Bearer %v", GetConfig().TwAPIToken)).
		Reply(200).
		Type("application/json").
		BodyString(mockedAccounts[2])

	api := NewTransferwiseAPI(GetConfig().TwHost, GetConfig().TwAPIToken)
	transfers, err := api.Transfers()
	require.NoError(t, err)
	require.Len(t, transfers, 2)

	transfer := transfers[0]
	require.Equal(t, int64(877677), transfer.ID)
	require.Equal(t, time.Date(2018, 8, 29, 17, 58, 24, 0, time.UTC), transfer.CreatedAt)
	require.Equal(t, "incoming_payment_waiting", transfer.Status)
	require.Equal(t, "Neo", transfer.RecipientName)
	require.Equal(t, 29894.57, transfer.SourceValue)
	require.Equal(t, "GBP", transfer.SourceCurrency)
	require.Equal(t, 33239.77, transfer.TargetValue)
	require.Equal(t, "EUR", transfer.TargetCurrency)
	// require.Equal(t, , transfer.Fee)
	require.Equal(t, 1.1119, transfer.ExchangeRate)

	transfer = transfers[1]
	require.Equal(t, int64(877677), transfer.ID)
	require.Equal(t, time.Date(2018, 8, 29, 17, 58, 24, 0, time.UTC), transfer.CreatedAt)
	require.Equal(t, "incoming_payment_waiting", transfer.Status)
	require.Equal(t, "Morpheus", transfer.RecipientName)
	require.Equal(t, 2989.45, transfer.SourceValue)
	require.Equal(t, "GBP", transfer.SourceCurrency)
	require.Equal(t, 3323.97, transfer.TargetValue)
	require.Equal(t, "EUR", transfer.TargetCurrency)
	// require.Equal(t, , transfer.Fee)
	require.Equal(t, 1.1119, transfer.ExchangeRate)
}
