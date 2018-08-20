package core_test

import (
	"testing"

	"github.com/endiangroup/transferwiser/core"
	coreMocks "github.com/endiangroup/transferwiser/core/mocks"
	kvMocks "github.com/endiangroup/transferwiser/keyvalue/mocks"
	"github.com/stretchr/testify/require"
)

func TestRedirectUrl(t *testing.T) {
	data := []struct {
		host       string
		clientID   string
		redirectCb string
		expected   string
	}{
		{
			host:       "test.com",
			clientID:   "foo/bar",
			redirectCb: "http://redirect.here/now",
			expected:   "https://test.com/oauth/authorize?response_type=code&client_id=foo%2Fbar&redirect_uri=http:%2F%2Fredirect.here%2Fnow",
		},
	}
	for _, test := range data {
		api := &coreMocks.TransferwiseAPI{}
		store := &kvMocks.Value{}
		service := core.NewTransferwiseService(api, test.clientID, test.host, test.redirectCb, store)
		require.Equal(t, test.expected, service.RedirectUrl())
	}
}

func TestUseAuthentication(t *testing.T) {
	token := "abcde1234"
	data := &core.RefreshTokenData{
		AccessToken: token,
	}
	api := &coreMocks.TransferwiseAPI{}
	store := &kvMocks.Value{}

	store.On("PutString", token).Return(nil)

	service := core.NewTransferwiseService(api, "foo/bar", "test.com", "http://redirect.now", store)
	require.False(t, service.IsAuthenticated())
	require.NoError(t, service.UseAuthentication(data))
	require.True(t, service.IsAuthenticated())
}
