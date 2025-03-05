package client

import (
	"encoding/json"
	"fmt"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
	"github.com/loungeup/go-loungeup/transport"
)

type segmentsClient struct{ baseClient *Client }

func (c *segmentsClient) BuildESQuery(
	selector *models.SegmentSelector,
	params *models.SearchCriterion,
) (*models.BuildSegmentESQueryResponse, error) {
	return transport.CallRESResult[*models.BuildSegmentESQueryResponse](
		c.baseClient.resClient,
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
