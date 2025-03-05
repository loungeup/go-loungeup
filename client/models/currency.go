package models

import "time"

type CurrencyRates struct {
	AUD float64 `json:"aud"`
	CAD float64 `json:"cad"`
	CHF float64 `json:"chf"`
	CNY float64 `json:"cny"`
	EUR float64 `json:"eur"`
	GBP float64 `json:"gbp"`
	JPY float64 `json:"jpy"`
	KRW float64 `json:"krw"`
	SGD float64 `json:"sgd"`
	USD float64 `json:"usd"`
}

type CurrencyRatesSelector struct {
	Base string
	Date time.Time
}

func (s CurrencyRatesSelector) RID() string { return "authority.currencies." + s.Base + ".rates" }

func (s CurrencyRatesSelector) EncodedQuery() string {
	return "date=" + s.Date.Format(time.DateOnly)
}
