package price

import (
	"log"
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

type PriceTestServer struct {
}

func (s *PriceTestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("asdfasd")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"data": {"priceUsd": "789.01"}}`))
}

func TestGetCurrentPrice_CoincapSuccess(t *testing.T) {
	server := httptest.NewServer(&PriceTestServer{})
	defer server.Close()

	ps := NewPriceService()

	// override urls to use mock test server
	for k := range ps.urls {
		ps.urls[k] = server.URL
	}

	price, err := ps.GetCurrentPrice(CrytpoAssetEthereum)
	assert.NoError(t, err)
	assert.Equal(t, 789.01, price)
}

func TestGetCurrentPrice_UnsupportedCryptoAsset(t *testing.T) {
	ps := NewPriceService()
	price, err := ps.GetCurrentPrice(CrytpoAssetUnsupportedCrypto)
	assert.Error(t, err)
	assert.Equal(t, 0.0, price)
}
