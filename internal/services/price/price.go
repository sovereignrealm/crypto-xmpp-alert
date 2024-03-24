package price

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type PriceResponse interface {
	getPrice(assetName CrytpoAsset) float64
}

type CoindeskResponse struct {
	BPI struct {
		USD struct {
			Rate float64 `json:"rate_float"`
		} `json:"USD"`
	} `json:"bpi"`
}

func (c *CoindeskResponse) getPrice(assetName CrytpoAsset) float64 {
	return c.BPI.USD.Rate
}

type CoinGeckoResponse struct {
	Cardano struct {
		USD float64 `json:"usd"`
	} `json:"cardano"`

	Polkadot struct {
		USD float64 `json:"usd"`
	} `json:"polkadot"`
}

func (c *CoinGeckoResponse) getPrice(assetName CrytpoAsset) float64 {
	switch assetName {
	case CrytpoAssetCardano:
		return c.Cardano.USD
	case CrytpoAssetPolkadot:
		return c.Polkadot.USD
	default:
		fmt.Printf("Warning: CoinGecko does not support %s.\n", assetName)
		return 0
	}
}

type CoinCapResponse struct {
	Data struct {
		PriceUsd string `json:"priceUsd"`
	} `json:"data"`
}

// type HTTPClient interface {
// 	Get(url string) (*http.Response, error)
// }

type PriceService struct {
	client *http.Client
	urls   map[string]string
}

func NewPriceService() PriceService {
	return PriceService{
		client: http.DefaultClient,
		urls: map[string]string{
			"Bitcoin":  "https://api.coindesk.com/v1/bpi/currentprice.json",
			"Ethereum": "https://api.coincap.io/v2/assets/ethereum",
			"Cardano":  "https://api.coingecko.com/api/v3/simple/price?ids=cardano&vs_currencies=usd",
			"Polkadot": "https://api.coingecko.com/api/v3/simple/price?ids=polkadot&vs_currencies=usd",
		},
	}
}

func (ps *PriceService) GetCurrentPrice(cryptoAsset CrytpoAsset) (float64, error) {
	url, found := ps.urls[string(cryptoAsset)]
	if !found {
		return 0, fmt.Errorf("unsupported crypto asset: %s", cryptoAsset)
	}

	resp, err := ps.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response PriceResponse

	switch cryptoAsset {
	case CrytpoAssetEthereum:
		var resp CoinCapResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return 0, err
		}
		priceUsd, err := strconv.ParseFloat(resp.Data.PriceUsd, 64)
		if err != nil {
			return 0, err
		}
		return priceUsd, nil
	case CrytpoAssetBitcoin:
		response = &CoindeskResponse{}
	case CrytpoAssetCardano, CrytpoAssetPolkadot:
		response = &CoinGeckoResponse{}
	default:
		return 0.0, fmt.Errorf("unsupported crypto asset: %s", cryptoAsset)
	}

	if err := json.Unmarshal(body, response); err != nil {
		return 0, err
	}
	return response.getPrice(cryptoAsset), nil
}
