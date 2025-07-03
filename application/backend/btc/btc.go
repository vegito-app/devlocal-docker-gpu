package btc

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type EstimatredPrice struct {
	Value       float64
	Cardinal    int
	TrustIndice float64
}

type BTC struct {
	stop               func()
	priceLock          sync.RWMutex
	EstimatedEuroPrice *EstimatredPrice
}

func NewBTC() *BTC {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	btc := &BTC{
		stop: func() {
			cancel()
			<-done
		},
	}
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Minute):
				btc.updatePrice(ctx)
			}
		}
	}()
	btc.updatePrice(ctx)
	return btc
}

func (btc *BTC) Close() {
	if btc != nil {
		btc.stop()
	}
	btc.stop = nil
}

func (btc *BTC) Price(resp http.ResponseWriter, req *http.Request) {
	btc.priceLock.RLock()
	defer btc.priceLock.RUnlock()
	if err := json.NewEncoder(resp).Encode(map[string]any{
		"btc_price_eur": btc.EstimatedEuroPrice,
	}); err != nil {
		log.Error().
			Err(err).
			Msg("backend btc price response")
	}
}

func (btc *BTC) setEstimatedEuroPrice(btcEuroPrice float64, trustIndice float64, valueCount int) {
	btc.priceLock.Lock()
	defer btc.priceLock.Unlock()
	btc.EstimatedEuroPrice = &EstimatredPrice{
		Value:       btcEuroPrice,
		Cardinal:    valueCount,
		TrustIndice: trustIndice,
	}
}

func (btc *BTC) updatePrice(ctx context.Context) {
	results := make(chan *apiPrice)
	go func() {
		defer close(results)

		var wg sync.WaitGroup
		defer wg.Wait()

		wg.Add(len(apis))
		for _, api := range apis {
			go func() {
				defer wg.Done()

				result, err := fetchPrice(ctx, &api)
				if err != nil {
					log.Error().
						Err(err).
						Str("api_name", api.Name).
						Msg("fetch api BTC price")
					return
				}
				results <- result
			}()
		}
	}()
	var rawPrices []*apiPrice
	for result := range results {
		rawPrices = append(rawPrices, result)
	}
	btcEuroPrice, filteredPrices := filterOutliers(rawPrices)
	trustedValuesCount := len(filteredPrices)
	trustIndice := float64(trustedValuesCount) / float64(len(rawPrices))
	if trustIndice < 0.5 {
		log.Error().
			Float64("btc_euro_price", btcEuroPrice).
			Int("used_api_number", len(rawPrices)).
			Int("filtered_api_number", trustedValuesCount).
			Float64("trust_indice", trustIndice).
			Msg("BTC euro price not updated")
	}
	btc.setEstimatedEuroPrice(btcEuroPrice, trustIndice, trustedValuesCount)
	log.Info().
		Float64("btc_euro_price", btcEuroPrice).
		Int("used_api_number", len(rawPrices)).
		Int("filtered_api_number", trustedValuesCount).
		Float64("trust_indice", trustIndice).
		Msg("updated BTC euro price")
}
