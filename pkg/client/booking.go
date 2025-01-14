package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/pkg/client/models"
)

type bookingsClient struct{ baseClient *Client }

func (c *bookingsClient) CountBookings(entityID uuid.UUID) (int64, error) {
	request, err := http.NewRequest(
		http.MethodPost,
		c.baseClient.httpAPIURL+"/entities/"+entityID.String()+"/bookings/count",
		http.NoBody,
	)
	if err != nil {
		return 0, fmt.Errorf("could not create request: %w", err)
	}

	response, err := c.baseClient.executeHTTPRequest(request)
	if err != nil {
		return 0, fmt.Errorf("could not send request: %w", err)
	}
	defer response.Body.Close()

	var result int64
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("could not decode response body: %w", err)
	}

	return result, nil
}

func (c *bookingsClient) ReadBookingIDs(selector *models.BookingIDsSelector) (models.ReadBookingIDsResponse, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		c.baseClient.httpAPIURL+"/entities/"+selector.EntityID.String()+"/booking-ids?"+selector.EncodedQuery(),
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	response, err := c.baseClient.executeHTTPRequest(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	defer response.Body.Close()

	result := models.ReadBookingIDsResponse{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w", err)
	}

	return result, nil
}

func (c *bookingsClient) ReadIndexableBookingByID(bookingID int) (*models.IndexableBookingResponse, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		c.baseClient.httpAPIURL+"/booking/"+strconv.Itoa(bookingID)+"/indexable",
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	response, err := c.baseClient.executeHTTPRequest(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	defer response.Body.Close()

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

func (c *bookingsClient) Search(entityID uuid.UUID, selector models.SearchBookingsRequest) (*models.SearchBookingsResponse, error) {
	body, err := json.Marshal(selector)
	if err != nil {
		return nil, fmt.Errorf("could not marshal payload: %w", err)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		c.baseClient.httpAPIURL+"/entities/"+entityID.String()+"/bookings/search",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	response, err := c.baseClient.executeHTTPRequest(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	defer response.Body.Close()

	result := &models.SearchBookingsResponse{}
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("could not decode response body: %w", err)
	}

	return result, nil
}
