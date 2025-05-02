package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jirenius/go-res/resprot"
	"github.com/loungeup/go-loungeup/client/models"
)

//go:generate mockgen -source booking.go -destination=./mocks/mock_booking.go -package=mocks

type BookingsManager interface {
	CountBookings(entityID uuid.UUID) (int64, error)
	ReadBookingIDs(selector *models.BookingIDsSelector) (models.ReadBookingIDsResponse, error)
	ReadIndexableBookingByID(bookingID int) (*models.IndexableBookingResponse, error)
	IndexBooking(request *models.IndexBookingRequest) error
	Search(entityID uuid.UUID, selector models.SearchBookingsRequest) (*models.SearchBookingsResponse, error)
}

type BookingsClient struct {
	base *BaseClient
}

func NewBookingsClient(base *BaseClient) *BookingsClient {
	return &BookingsClient{
		base: base,
	}
}

func (c *BookingsClient) CountBookings(entityID uuid.UUID) (int64, error) {
	request, err := http.NewRequest(
		http.MethodPost,
		c.base.httpAPIURL+"/entities/"+entityID.String()+"/bookings/count",
		http.NoBody,
	)
	if err != nil {
		return 0, fmt.Errorf("could not create request: %w", err)
	}

	response, err := c.base.ExecuteHTTPRequest(request)
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

func (c *BookingsClient) ReadBookingIDs(selector *models.BookingIDsSelector) (models.ReadBookingIDsResponse, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		c.base.httpAPIURL+"/entities/"+selector.EntityID.String()+"/booking-ids?"+selector.EncodedQuery(),
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	response, err := c.base.ExecuteHTTPRequest(request)
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

func (c *BookingsClient) ReadIndexableBookingByID(bookingID int) (*models.IndexableBookingResponse, error) {
	request, err := http.NewRequest(
		http.MethodGet,
		c.base.httpAPIURL+"/booking/"+strconv.Itoa(bookingID)+"/indexable",
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	response, err := c.base.ExecuteHTTPRequest(request)
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

func (c *BookingsClient) IndexBooking(request *models.IndexBookingRequest) error {
	if response := c.base.resClient.Request("call."+request.RID(), resprot.Request{
		Params: request,
	}); response.HasError() {
		return response.Error
	}

	return nil
}

func (c *BookingsClient) Search(entityID uuid.UUID, selector models.SearchBookingsRequest) (*models.SearchBookingsResponse, error) {
	body, err := json.Marshal(selector)
	if err != nil {
		return nil, fmt.Errorf("could not marshal payload: %w", err)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		c.base.httpAPIURL+"/entities/"+entityID.String()+"/bookings/search",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	response, err := c.base.ExecuteHTTPRequest(request)
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
