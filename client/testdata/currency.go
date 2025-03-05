//nolint:gomnd,mnd
package testdata

import (
	"time"

	"github.com/loungeup/go-loungeup/client/models"
)

var (
	CurrencyRates = &models.CurrencyRates{
		AUD: 1.642043,
		CAD: 1.512434,
		CHF: 0.939101,
		CNY: 7.888231,
		EUR: 1,
		GBP: 0.842507,
		JPY: 157.736305,
		KRW: 1472.87297,
		SGD: 1.440295,
		USD: 1.112837,
	}

	CurrencyRatesModel = `{
		"aud": 1.642043,
		"cad": 1.512434,
		"chf": 0.939101,
		"cny": 7.888231,
		"eur": 1,
		"gbp": 0.842507,
		"jpy": 157.736305,
		"krw": 1472.87297,
		"sgd": 1.440295,
		"usd": 1.112837
	}`

	CurrencyRatesSelector = models.CurrencyRatesSelector{
		Base: "eur",
		Date: time.Date(2024, time.September, 18, 0, 0, 0, 0, time.UTC),
	}
)
