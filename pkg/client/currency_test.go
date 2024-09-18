package client

import (
	"strings"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/testdata"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/loungeup/go-loungeup/pkg/transporttest"
	"github.com/stretchr/testify/assert"
)

func TestReadCurrencyRates(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.CurrencyRatesSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.CurrencyRatesModel)
				default:
					return transporttest.NewRESModelResponse(`{}`)
				}
			},
		},
	}).Internal.Currency.ReadCurrencyRates(&testdata.CurrencyRatesSelector)
	assert.NoError(t, err)
	assert.Equal(t, testdata.CurrencyRates, got)
}
