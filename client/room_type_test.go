package client

import (
	"strings"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/client/testdata"
	"github.com/loungeup/go-loungeup/transport"
	"github.com/loungeup/go-loungeup/transporttest"
	"github.com/stretchr/testify/assert"
)

func TestReadRoomTypes(t *testing.T) {
	got, err := NewWithTransport(
		&transport.Transport{
			RESClient: &transporttest.RESClientMock{
				RequestFunc: func(resourceID string, request resprot.Request) resprot.Response {
					switch {
					case strings.HasSuffix(resourceID, testdata.RoomTypesSelector.RID()):
						return transporttest.NewRESCollectionResponse(testdata.RoomTypeCollection)
					case strings.HasSuffix(resourceID, testdata.RoomTypeSelector.RID()):
						return transporttest.NewRESModelResponse(testdata.RoomTypeModel)
					default:
						return transporttest.NewRESModelResponse(`{}`)
					}
				},
			},
		},
	).RoomTypes.ReadRoomTypes(testdata.RoomTypesSelector)
	assert.NoError(t, err)
	assert.Equal(t, []*models.RoomType{testdata.RoomType}, got)
}
