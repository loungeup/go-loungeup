package client

import (
	"time"

	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/resmodels"
	"github.com/loungeup/go-loungeup/transport"
)

type ComputedAttrsManager interface {
	ReadOne(selector *ComputedAttrSelector) (*resmodels.ComputedAttr, error)
}

type ComputedAttrsClient struct {
	base *BaseClient
}

func NewComputedAttrsClient(base *BaseClient) *ComputedAttrsClient {
	return &ComputedAttrsClient{
		base: base,
	}
}

func (c *ComputedAttrsClient) ReadOne(selector *ComputedAttrSelector) (*resmodels.ComputedAttr, error) {
	cacheKey := selector.rid()

	if cachedResult, ok := c.base.ReadCache(cacheKey).(*resmodels.ComputedAttr); ok {
		return cachedResult, nil
	}

	result, err := transport.GetRESModel[*resmodels.ComputedAttr](
		c.base.resClient,
		selector.rid(),
		resprot.Request{},
	)
	if err != nil {
		return nil, err
	}

	c.base.WriteCacheWithDuration(cacheKey, result, time.Minute)

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
