package btc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

// API configuration
type APIConfig struct {
	Name    string
	URL     string
	ParseFn func([]byte) (float64, error)
}

// List of APIs
var apis = []APIConfig{
	{
		Name: "CoinGecko",
		URL:  "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=eur",
		ParseFn: func(body []byte) (float64, error) {
			var data map[string]map[string]float64
			if err := json.Unmarshal(body, &data); err != nil {
				return 0, fmt.Errorf("parsing CoinGecko price: %w", err)
			}
			return data["bitcoin"]["eur"], nil
		},
	},
	// Add more APIs (e.g., CoinMarketCap)
	{
		Name: "CryptoCompare",
		URL:  "https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=EUR",
		ParseFn: func(body []byte) (float64, error) {
			var data map[string]float64
			if err := json.Unmarshal(body, &data); err != nil {
				return 0, fmt.Errorf("parsing CryptoCompare price: %w", err)
			}
			return data["EUR"], nil
		},
	},
	// Add more APIs (e.g., Binance)
	{
		Name: "Binance",
		URL:  "https://api.binance.com/api/v3/ticker/price?symbol=BTCEUR",
		ParseFn: func(body []byte) (float64, error) {
			var data map[string]string
			if err := json.Unmarshal(body, &data); err != nil {
				return 0, fmt.Errorf("parsing Binance price: %w", err)
			}
			price, err := strconv.ParseFloat(data["price"], 64)
			if err != nil {
				return 0, fmt.Errorf("convert binance euro price from string: %w", err)
			}
			return price, nil
		},
	},
	// Add more APIs (e.g., Binance)
	{
		Name: "Kraken",
		URL:  "https://api.kraken.com/0/public/Ticker?pair=BTCEUR",
		ParseFn: func(body []byte) (float64, error) {
			var data struct {
				Result map[string]struct {
					Ask []string `json:"a"`
				} `json:"result"`
			}
			if err := json.Unmarshal(body, &data); err != nil {
				return 0, fmt.Errorf("parsing Kraken price: %w", err)
			}
			for _, ticker := range data.Result {
				price, err := strconv.ParseFloat(ticker.Ask[0], 64)
				if err != nil {
					return 0, fmt.Errorf("parsing Kraken price: %w", err)
				}
				return price, nil
			}
			return 0, fmt.Errorf("no Kraken BTC/EUR price found")
		},
	},
	// Add more APIs (e.g., Bitfinex)
	{
		Name: "Bitfinex",
		URL:  "https://api.binance.com/api/v3/ticker/price?symbol=BTCEUR",
		ParseFn: func(body []byte) (float64, error) {
			var data map[string]any
			if err := json.Unmarshal(body, &data); err != nil {
				return 0, fmt.Errorf("parsing Bitfinex price: %w", err)
			}
			priceStr, ok := data["price"].(string)
			if !ok {
				return 0, fmt.Errorf("parsing Bitfinex price")
			}
			price, err := strconv.ParseFloat(priceStr, 64)
			if err != nil {
				return 0, fmt.Errorf("parse Bitfinex BTC/EUR price from string value: %w", err)
			}
			return price, nil
		},
	},
	// Add more APIs (e.g., Coinbase)
	{
		Name: "Coinbase",
		URL:  "https://api.coinbase.com/v2/prices/BTC-EUR/spot",
		ParseFn: func(body []byte) (float64, error) {
			var data struct {
				Data struct {
					Amount string `json:"amount"`
				} `json:"data"`
			}
			if err := json.Unmarshal(body, &data); err != nil {
				return 0, fmt.Errorf("parsing Coinbase price: %w", err)
			}
			price, err := strconv.ParseFloat(data.Data.Amount, 64)
			if err != nil {
				return 0, fmt.Errorf("parsing Coinbase price: %w", err)
			}
			return price, nil
		},
	},
	// Add more APIs (e.g., OKX)
	{
		Name: "okx",
		URL:  "https://www.okx.com/api/v5/market/ticker?instId=BTC-EUR",
		ParseFn: func(body []byte) (float64, error) {
			var data struct {
				Data []struct {
					Last string `json:"last"`
				} `json:"data"`
			}
			if err := json.Unmarshal(body, &data); err != nil {
				return 0, fmt.Errorf("parsing okx price: %w", err)
			}
			if len(data.Data) > 0 {
				price, err := strconv.ParseFloat(data.Data[0].Last, 64)
				if err != nil {
					return 0, fmt.Errorf("parsing okx price: %v", err)
				}
				return price, nil
			}

			return 0, fmt.Errorf("no okx BTC/EUR price found")
		},
	},
}

type apiPrice struct {
	api   *APIConfig
	value float64
}

// Function to fetch price from one API
func fetchPrice(ctx context.Context, api *APIConfig) (results *apiPrice, errors error) {
	req, err := http.NewRequestWithContext(ctx, "GET", api.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("new API %s, fetch request: %w", api.Name, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API %s failed: %w", api.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API %s returned status: %d", api.Name, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from %s: %w", api.Name, err)
	}

	price, err := api.ParseFn(body)
	if err != nil {
		return nil, fmt.Errorf("parsing data from %s: %w", api.Name, err)
	}

	return &apiPrice{
		api:   api,
		value: price,
	}, nil
}

func filterOutliers(prices []*apiPrice) (float64, []*apiPrice) {
	if len(prices) == 0 {
		return 0, nil
	}

	// Calculer la moyenne
	sum := 0.0
	for _, price := range prices {
		sum += price.value
	}
	mean := sum / float64(len(prices))

	// Calculer l'écart-type
	sumSquareDiffs := 0.0
	for _, price := range prices {
		diff := price.value - mean
		sumSquareDiffs += diff * diff
	}
	stdDev := math.Sqrt(sumSquareDiffs / float64(len(prices)))

	// Filtrer les valeurs hors des 2 écarts-types
	filteredPrices := []*apiPrice{}
	for _, apiPrice := range prices {
		price, apiName := apiPrice.value, apiPrice.api.Name
		if math.Abs(price-mean) > 2*stdDev {
			log.Warn().
				Str("api_name", apiName).
				Float64("btc_euro_price", price).
				Msg("discarding BTC euro price")
			continue
		}
		filteredPrices = append(filteredPrices, apiPrice)
	}

	// Recalculer la moyenne
	if len(filteredPrices) == 0 {
		return mean, filteredPrices
	}
	sum = 0.0
	for _, price := range filteredPrices {
		sum += price.value
	}
	return sum / float64(len(filteredPrices)), filteredPrices
}
