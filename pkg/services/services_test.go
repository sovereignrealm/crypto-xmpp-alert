package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoindeskResponse_getPrice(t *testing.T) {
	coindeskResponse := &CoindeskResponse{
		BPI: struct {
			USD struct {
				Rate float64 `json:"rate_float"`
			} `json:"USD"`
		}{
			USD: struct {
				Rate float64 `json:"rate_float"`
			}{Rate: 123.45},
		},
	}

	price := coindeskResponse.getPrice("Bitcoin")

	assert.Equal(t, 123.45, price)
}

func TestCoinGeckoResponse_getPrice_Cardano(t *testing.T) {
	coinGeckoResponse := &CoinGeckoResponse{
		Cardano: struct {
			USD float64 `json:"usd"`
		}{USD: 1.23},
	}

	price := coinGeckoResponse.getPrice("Cardano")

	assert.Equal(t, 1.23, price)
}

func TestCoinGeckoResponse_getPrice_Polkadot(t *testing.T) {
	coinGeckoResponse := &CoinGeckoResponse{
		Polkadot: struct {
			USD float64 `json:"usd"`
		}{USD: 45.67},
	}

	price := coinGeckoResponse.getPrice("Polkadot")

	assert.Equal(t, 45.67, price)
}

func TestGetCurrentPrice_CoincapSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data": {"priceUsd": "789.01"}}`))
	}))
	defer server.Close()

	mockHTTPClient := MockHTTPClient{
		RespBody: []byte(`{"data": {"PriceUsd": "789.01"}}`),
	}
	price, err := GetCurrentPrice(&mockHTTPClient, "Ethereum")
	assert.NoError(t, err)
	assert.Equal(t, 789.01, price)
}

func TestGetCurrentPrice_UnsupportedCryptoAsset(t *testing.T) {
	mockHTTPClient := MockHTTPClient{
		RespBody: []byte(`{"data": {"PriceUsd": "200"}}`),
	}
	price, err := GetCurrentPrice(&mockHTTPClient, "UnsupportedCrypto")
	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
}
