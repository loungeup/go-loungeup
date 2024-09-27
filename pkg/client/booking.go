package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
)

type bookingsClient struct{ baseClient *Client }

func (c *bookingsClient) ReadIndexableBookingByID(bookingID int) (*models.IndexableBookingResponse, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		c.baseClient.httpAPIURL+"/booking/"+strconv.Itoa(bookingID)+"/indexable",
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+c.baseClient.httpAPIKey)

	response, err := c.baseClient.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	result := &models.IndexableBookingResponse{}
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w", err)
	}

	return result, nil
}

func (c *bookingsClient) IndexBooking(request *models.IndexBookingRequest) error {
	if response := c.baseClient.resClient.Request("call."+request.RID(), resprot.Request{
		Params: request,
	}); response.HasError() {
		return response.Error
	}

	return nil
}
