package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type CoindeskResponse struct {
	BPI struct {
		USD struct {
			Rate float64 `json:"rate_float"`
		} `json:"USD"`
	} `json:"bpi"`
}

type CoinGeckoResponse struct {
	Cardano struct {
		USD float64 `json:"usd"`
	} `json:"cardano"`

	Polkadot struct {
		USD float64 `json:"usd"`
	} `json:"polkadot"`
}

type CoinCapResponse struct {
	Data struct {
		PriceUsd string `json:"priceUsd"`
	} `json:"data"`
}

type PriceResponse interface {
	getPrice(assetName string) float64
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

func (c *CoindeskResponse) getPrice(assetName string) float64 {
	return c.BPI.USD.Rate
}

func (c *CoinGeckoResponse) getPrice(assetName string) float64 {
	switch assetName {
	case "Cardano":
		return c.Cardano.USD
	case "Polkadot":
		return c.Polkadot.USD
	default:
		fmt.Printf("Warning: CoinGecko does not support %s.\n", assetName)
		return 0
	}
}

func GetCurrentPrice(client HTTPClient, cryptoAsset string) (float64, error) {
	urls := map[string]string{
		"Bitcoin":  "https://api.coindesk.com/v1/bpi/currentprice.json",
		"Ethereum": "https://api.coincap.io/v2/assets/ethereum",
		"Cardano":  "https://api.coingecko.com/api/v3/simple/price?ids=cardano&vs_currencies=usd",
		"Polkadot": "https://api.coingecko.com/api/v3/simple/price?ids=polkadot&vs_currencies=usd",
	}

	url, ok := urls[cryptoAsset]
	if !ok {
		return 0, fmt.Errorf("unsupported crypto asset: %s", cryptoAsset)
	}

	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if cryptoAsset == "Ethereum" {
		var resp CoinCapResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return 0, err
		}
		priceUsd, err := strconv.ParseFloat(resp.Data.PriceUsd, 64)
		if err != nil {
			return 0, err
		}
		return priceUsd, nil
	}

	var response PriceResponse

	switch cryptoAsset {
	case "Bitcoin":
		response = &CoindeskResponse{}
	case "Cardano", "Polkadot":
		response = &CoinGeckoResponse{}
	default:
		return 0.0, fmt.Errorf("unsupported crypto asset: %s", cryptoAsset)
	}

	if err := json.Unmarshal(body, response); err != nil {
		return 0, err
	}
	return response.getPrice(cryptoAsset), nil
}
