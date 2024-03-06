package client

import (
	"strings"
	"testing"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/client/testdata"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/stretchr/testify/assert"
)

func TestReadRoomTypes(t *testing.T) {
	got, err := NewWithTransport(&transport.Transport{
		RESClient: &resClientMock{
			requestFunc: func(resourceID string, request resprot.Request) resprot.Response {
				switch {
				case strings.HasSuffix(resourceID, testdata.RoomTypesSelector.RID()):
					return newCollectionResponse(testdata.RoomTypeCollection)
				case strings.HasSuffix(resourceID, testdata.RoomTypeSelector.RID()):
					return newModelResponse(testdata.RoomTypeModel)
				default:
					return newModelResponse(`{}`)
				}
			},
		},
	}).Internal.RoomTypes.ReadRoomTypes(testdata.RoomTypesSelector)
	assert.NoError(t, err)
	assert.Equal(t, []models.RoomType{testdata.RoomType}, got)
}
