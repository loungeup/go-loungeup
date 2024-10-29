package client

import (
	"encoding/json"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
)

type segmentsClient struct{ baseClient *Client }

func (c *segmentsClient) BuildESQuery(
	selector *models.SegmentSelector,
	query *models.SegmentQuery,
) (*models.BuildESQueryResponse, error) {
	return transport.CallRESResult[*models.BuildESQueryResponse](
		c.baseClient.resClient,
		selector.RID(),
		resprot.Request{
			Params: query,
			Token:  json.RawMessage(`{"agentRoles": ["service"]}`),
		},
	)
}
