package client

import (
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

type roomTypesClient struct{ baseClient *Client }

func (c *roomTypesClient) ReadRoomTypes(selector models.RoomTypesSelector) ([]models.RoomType, error) {
	if cachedResult, ok := c.baseClient.eventuallyReadCache(selector.RID()).([]models.RoomType); ok {
		return cachedResult, nil
	}

	references, err := transport.GetRESCollection[res.Ref](c.baseClient.resClient, selector.RID(), resprot.Request{})
	if err != nil {
		return nil, err
	}

	result := []models.RoomType{}

	for _, reference := range references {
		relatedRoomType, err := transport.GetRESModel[models.RoomType](c.baseClient.resClient, string(reference))
		if err != nil {
			return nil, err
		}

		result = append(result, relatedRoomType)
	}

	defer c.baseClient.eventuallyWriteCache(selector.RID(), result)

	return result, nil
}
