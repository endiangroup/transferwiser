package web

// type testData struct {
// 	transferwiseAPI *coreMocks.TransferwiseAPI
// 	webApp          *echo.Echo
// }

// func getTestData(t *testing.T) *testData {
// 	transferwiseAPI := &coreMocks.TransferwiseAPI{}

// 	logger, err := zap.NewDevelopment()
// 	require.NoError(t, err)
// 	webServer := NewServer(logger, transferwiseAPI)

// 	webApp := webServer.MainHandler()
// 	return &testData{
// 		transferwiseAPI: transferwiseAPI,
// 		webApp:          webApp,
// 	}
// }

// func TestTransferwiseTransfers_RespondsWithACsv(t *testing.T) {
// 	data := getTestData(t)

// 	transfers := []*core.Transfer{
// 		{
// 			ID:             1,
// 			CreatedAt:      time.Now(),
// 			Status:         "incoming_payment_waiting",
// 			RecipientName:  "Bill Gates",
// 			SourceValue:    10.000,
// 			SourceCurrency: "GPB",
// 			TargetValue:    12.3456,
// 			TargetCurrency: "EUR",
// 			Fee:            1.234,
// 			ExchangeRate:   1.2345,
// 		},
// 	}
// 	data.transferwiseAPI.On("Transfers").Return(transfers, nil)

// 	req, err := http.NewRequest(echo.GET, "/transfers.csv", nil)
// 	require.NoError(t, err)
// 	req.Header.Set("DN", "EMAIL=admin@endian.io,CN=admin,C=UK")
// 	req.Header.Set("VERIFIED", "SUCCESS")

// 	resRec := httptest.NewRecorder()
// 	data.webApp.ServeHTTP(resRec, req)

// 	require.Equal(t, 200, resRec.Code)
// 	require.Equal(t, "text/csv", resRec.Header().Get(echo.HeaderContentType))

// 	expected, err := gocsv.MarshalString(&transfers)
// 	require.NoError(t, err)
// 	require.Equal(t, expected, resRec.Body.String())
// }

// func TestTransferwiseTransfers_RequiresCN(t *testing.T) {
// 	data := getTestData(t)

// 	req, err := http.NewRequest(echo.GET, "/transfers.csv", nil)
// 	require.NoError(t, err)
// 	req.Header.Set("DN", "EMAIL=admin@endian.io,C=UK")
// 	req.Header.Set("VERIFIED", "SUCCESS")

// 	resRec := httptest.NewRecorder()
// 	data.webApp.ServeHTTP(resRec, req)

// 	require.Equal(t, 401, resRec.Code)
// }

// func TestTransferwiseTransfers_RequiresVerifiedCertificate(t *testing.T) {
// 	data := getTestData(t)

// 	req, err := http.NewRequest(echo.GET, "/transfers.csv", nil)
// 	require.NoError(t, err)
// 	req.Header.Set("DN", "EMAIL=admin@endian.io,CN=admin,C=UK")
// 	req.Header.Set("VERIFIED", "NONE")

// 	resRec := httptest.NewRecorder()
// 	data.webApp.ServeHTTP(resRec, req)

// 	require.Equal(t, 401, resRec.Code)
// }
