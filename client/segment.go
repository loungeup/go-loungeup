package client

import (
	"encoding/json"
	"fmt"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/transport"
)

//go:generate mockgen -source segment.go -destination=./mocks/mock_segment.go -package=mocks

type SegmentsManager interface {
	BuildESQuery(
		selector *models.SegmentSelector,
		params *models.SearchCriterion,
	) (*models.BuildSegmentESQueryResponse, error)
}

type SegmentsClient struct {
	base *BaseClient
}

func NewSegmentsClient(base *BaseClient) *SegmentsClient {
	return &SegmentsClient{
		base: base,
	}
}

func (c *SegmentsClient) BuildESQuery(
	selector *models.SegmentSelector,
	params *models.SearchCriterion,
) (*models.BuildSegmentESQueryResponse, error) {
	return transport.CallRESResult[*models.BuildSegmentESQueryResponse](
		c.base.resClient,
		fmt.Sprintf("guestprofile.entities.%s.segments.%s.build-elasticsearch-query",
			selector.EntityID.String(),
			selector.SegmentID.String(),
		),
		resprot.Request{
			Params: params,
			Token:  json.RawMessage(`{"agentRoles": ["service"]}`),
		},
	)
}
