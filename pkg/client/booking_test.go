package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/loungeup/go-loungeup/pkg/client/models"
	"github.com/loungeup/go-loungeup/pkg/transport"
	"github.com/stretchr/testify/require"
)

func TestCountBookings(t *testing.T) {
	entityID := uuid.New()

	tests := map[string]struct {
		handlerFunc http.HandlerFunc
		assertFunc  func(t *testing.T, got int64, err error)
	}{
		"simple": {
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "/entities/"+entityID.String()+"/bookings/count", r.URL.Path)
				require.Equal(t, "Bearer key", r.Header.Get("Authorization"))

				w.Write(json.RawMessage(`10`))
			},
			assertFunc: func(t *testing.T, got int64, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(10), got)
			},
		},
		"invalid status code": {
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			assertFunc: func(t *testing.T, got int64, err error) {
				require.Error(t, err)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewWithTransport(
				transport.New(),
				WithHTTPAPIKey("key"),
				WithHTTPAPIURL(httptest.NewServer(tt.handlerFunc).URL),
			).Internal.Bookings.CountBookings(entityID)
			tt.assertFunc(t, got, err)
		})
	}
}

func TestReadBookingIDs(t *testing.T) {
	tests := map[string]struct {
		selector    *models.BookingIDsSelector
		handlerFunc http.HandlerFunc
		assertFunc  func(t *testing.T, got models.ReadBookingIDsResponse, err error)
	}{
		"simple": {
			selector: &models.BookingIDsSelector{
				EntityID:    uuid.MustParse("0ce1e959-57e3-4ebb-b8d9-6126a05afee2"),
				Size:        20,
				LastGuestID: uuid.MustParse("ccb8622e-e70f-41ef-bb0c-94b6dcfeca68"),
			},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "/entities/0ce1e959-57e3-4ebb-b8d9-6126a05afee2/booking-ids", r.URL.Path)
				require.Equal(t, "20", r.URL.Query().Get("size"))
				require.Equal(t, "ccb8622e-e70f-41ef-bb0c-94b6dcfeca68", r.URL.Query().Get("lastGuestId"))
				require.Equal(t, "Bearer key", r.Header.Get("Authorization"))

				w.Write(json.RawMessage(`[
					{
						"id": 22206329,
						"guestId": "02c774a6-71e8-94db-4961-28ce29f0e729"
					},
					{
						"id": 22206334,
						"guestId": "2da87948-af4f-f0a7-8d39-05a7f464b946"
					},
					{
						"id": 22206332,
						"guestId": "3412d5d9-c344-30c3-d0e8-6c8d20027eec"
					},
					{
						"id": 22206330,
						"guestId": "39d209ea-fc4e-8d31-8146-2fd456943874"
					}
				]`))
			},
			assertFunc: func(t *testing.T, got models.ReadBookingIDsResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, models.ReadBookingIDsResponse{
					{
						ID:      22206329,
						GuestID: uuid.MustParse("02c774a6-71e8-94db-4961-28ce29f0e729"),
					},
					{
						ID:      22206334,
						GuestID: uuid.MustParse("2da87948-af4f-f0a7-8d39-05a7f464b946"),
					},
					{
						ID:      22206332,
						GuestID: uuid.MustParse("3412d5d9-c344-30c3-d0e8-6c8d20027eec"),
					},
					{
						ID:      22206330,
						GuestID: uuid.MustParse("39d209ea-fc4e-8d31-8146-2fd456943874"),
					},
				}, got)
			},
		},
		"invalid status code": {
			selector:    &models.BookingIDsSelector{EntityID: uuid.MustParse("0ce1e959-57e3-4ebb-b8d9-6126a05afee2")},
			handlerFunc: func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) },
			assertFunc:  func(t *testing.T, _ models.ReadBookingIDsResponse, err error) { require.Error(t, err) },
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewWithTransport(
				transport.New(),
				WithHTTPAPIKey("key"),
				WithHTTPAPIURL(httptest.NewServer(tt.handlerFunc).URL),
			).Internal.Bookings.ReadBookingIDs(tt.selector)
			tt.assertFunc(t, got, err)
		})
	}
}
