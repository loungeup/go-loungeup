package client

import (
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/transport"
)

//go:generate mockgen -source room_type.go -destination=./mocks/mock_room_type.go -package=mocks

type RoomTypesManager interface {
	ReadRoomTypes(selector *models.RoomTypesSelector) ([]*models.RoomType, error)
}

type RoomTypesClient struct {
	base *BaseClient
}

func NewRoomTypesClient(base *BaseClient) *RoomTypesClient {
	return &RoomTypesClient{
		base: base,
	}
}

func (c *RoomTypesClient) ReadRoomTypes(selector *models.RoomTypesSelector) ([]*models.RoomType, error) {
	if cachedResult, ok := c.base.ReadCache(selector.RID()).([]*models.RoomType); ok {
		return cachedResult, nil
	}

	references, err := transport.GetRESCollection[res.Ref](c.base.resClient, selector.RID(), resprot.Request{})
	if err != nil {
		return nil, err
	}

	result := []*models.RoomType{}

	for _, reference := range references {
		relatedRoomType, err := transport.GetRESModel[*models.RoomType](
			c.base.resClient,
			string(reference),
			resprot.Request{},
		)
		if err != nil {
			return nil, err
		}

		result = append(result, relatedRoomType)
	}

	defer c.base.WriteCache(selector.RID(), result)

	return result, nil
}
