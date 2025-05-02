package client

import (
	"strings"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/testdata"
	"github.com/loungeup/go-loungeup/transport"
	"github.com/loungeup/go-loungeup/transporttest"
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
	}, nil).Currency.ReadCurrencyRates(&testdata.CurrencyRatesSelector)
	assert.NoError(t, err)
	assert.Equal(t, testdata.CurrencyRates, got)
}
