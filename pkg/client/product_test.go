package client

import (
	"strings"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/client/testdata"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/loungeup/go-loungeup/pkg/transporttest"
	"github.com/stretchr/testify/assert"
)

func TestReadProducts(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &transporttest.RESClientMock{
			RequestFunc: func(resourceID string, _ resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.ProductsSelector.RID()):
					return transporttest.NewRESCollectionResponse(testdata.ProductsCollection)
				case strings.HasSuffix(resourceID, testdata.ProductSelector.RID()):
					return transporttest.NewRESModelResponse(testdata.ProductModel)
				default:
					return transporttest.NewRESModelResponse(`{}`)
				}
			},
		},
	}).Internal.Products.ReadProducts(testdata.ProductsSelector)
	assert.NoError(t, err)
	assert.Equal(t, []*models.Product{testdata.Product}, got)
}
