package client

import (
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/resmodels"
	"github.com/loungeup/go-loungeup/transport"
)

type computedAttrsClient struct{ baseClient *Client }

func (c *computedAttrsClient) ReadOne(selector *ComputedAttrSelector) (*resmodels.ComputedAttr, error) {
	cacheKey := selector.rid()

	if cachedResult, ok := c.baseClient.eventuallyReadCache(cacheKey).(*resmodels.ComputedAttr); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*resmodels.ComputedAttr](
		c.baseClient.resClient,
		selector.rid(),
		resprot.Request{},
	)
	if err != nil {
		return nil, err
	}

	defer c.baseClient.eventuallyWriteCacheWithDuration(cacheKey, result, time.Minute)

	return result, nil
}

type ComputedAttrSelector struct {
	AttrID   uuid.UUID
	EntityID uuid.UUID
}

func (s *ComputedAttrSelector) rid() string {
	if s.EntityID == uuid.Nil {
		return "guestprofile.computed-attributes." + s.AttrID.String()
	}

	return "guestprofile.entities." + s.EntityID.String() + ".computed-attributes." + s.AttrID.String()
}
