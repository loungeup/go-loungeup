package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const aDay = 24 * time.Hour

type CustomField struct {
	Value string `json:"value"`
	From  string `json:"from"`
}

type CustomFields map[string]CustomField

type Booking struct {
	ID                 int          `json:"id"`
	EntityID           uuid.UUID    `json:"entityId"`
	GuestID            uuid.UUID    `json:"guestId"`
	Arrival            time.Time    `json:"arrival"`
	Departure          time.Time    `json:"departure"`
	Room               string       `json:"room"`
	RoomType           string       `json:"roomType"`
	PMSBookingID       string       `json:"pmsBookingId"`
	PMSBookingParentID string       `json:"pmsBookingParentId"`
	BookingDate        time.Time    `json:"bookingDate"`
	PaxAdults          string       `json:"paxAdults,omitempty"`
	PaxChildren        string       `json:"paxChildren,omitempty"`
	CustomFields       CustomFields `json:"customFields"`
}

func (booking *Booking) InStayDates() []time.Time {
	result := []time.Time{
		booking.arrivalDay(),
	}
	for departure := booking.arrivalDay().Add(aDay); departure.Before(booking.departureDay()); departure = departure.Add(aDay) {
		result = append(result, departure)
	}

	return result
}

func (booking *Booking) arrivalDay() time.Time { return booking.Arrival.Truncate(aDay) }

func (b *Booking) departureDay() time.Time {
	if b.Departure.IsZero() {
		return b.arrivalDay()
	}

	return b.Departure.Truncate(aDay)
}

type BookingSelector struct {
	BookingID int
	EntityID  uuid.UUID
}

func (s *BookingSelector) RID() string {
	return "proxy-db.entities." + s.EntityID.String() + ".bookings." + strconv.Itoa(s.BookingID)
}

type BookingSelectorById struct {
	BookingID int
}

func (selector *BookingSelectorById) RID() string {
	return "proxy-db.bookings." + strconv.Itoa(selector.BookingID)
}

func (c CustomFields) Get(key string) string {
	if field, ok := c[key]; ok {
		return field.Value
	}

	return ""
}

type IndexableBookingResponse struct {
	Booking RawIndexableBooking `json:"booking"`
}

type RawIndexableBooking struct {
	ID        int
	EntityID  uuid.UUID
	Departure time.Time

	Full json.RawMessage
}

var _ json.Marshaler = (*RawIndexableBooking)(nil)

func (b *RawIndexableBooking) MarshalJSON() ([]byte, error) {
	return b.Full, nil
}

var _ json.Unmarshaler = (*RawIndexableBooking)(nil)

func (b *RawIndexableBooking) UnmarshalJSON(data []byte) error {
	minimalBooking := &struct {
		ID        int       `json:"id"`
		EntityID  uuid.UUID `json:"entityId"`
		Departure time.Time `json:"departure"`
	}{}
	if err := json.Unmarshal(data, minimalBooking); err != nil {
		return fmt.Errorf("could not decode minimal booking: %w", err)
	}

	b.ID = minimalBooking.ID
	b.EntityID = minimalBooking.EntityID
	b.Departure = minimalBooking.Departure
	b.Full = data

	return nil
}

type IndexBookingRequest struct {
	Booking             RawIndexableBooking `json:"booking"`
	CampaignStats       json.RawMessage     `json:"campaignStats,omitempty"`
	ComputeAggregations bool                `json:"computeAggregations,omitempty"`
	ReindexGuest        *bool               `json:"reindexGuest,omitempty"`
	SurveyAnswers       json.RawMessage     `json:"surveyAnswers,omitempty"`
	UpdateGuestExtra    *bool               `json:"updateGuestExtra,omitempty"`
	UserDevice          json.RawMessage     `json:"userDevice,omitempty"`
}

func (r *IndexBookingRequest) RID() string {
	return strings.Join([]string{
		"indexer",
		"entities",
		r.Booking.EntityID.String(),
		"guest-bookings",
		strconv.Itoa(r.Booking.ID),
		"index",
	}, ".")
}

type BookingIDsSelector struct {
	EntityID    uuid.UUID
	Size        int
	LastGuestID uuid.UUID
}

func (s *BookingIDsSelector) EncodedQuery() string {
	result := "entityId=" + s.EntityID.String()

	if s.LastGuestID != uuid.Nil {
		result += "&lastGuestId=" + s.LastGuestID.String()
	}

	if s.Size > 0 {
		result += "&size=" + strconv.Itoa(s.Size)
	}

	return result
}

type ReadBookingIDsResponse []struct {
	ID      int       `json:"id"`
	GuestID uuid.UUID `json:"guestId"`
}

func (s ReadBookingIDsResponse) BookingIDs() []int {
	result := []int{}
	for _, bookingID := range s {
		result = append(result, bookingID.ID)
	}

	return result
}

func (s ReadBookingIDsResponse) LastGuestID() uuid.UUID {
	if len(s) == 0 {
		return uuid.Nil
	}

	return s[len(s)-1].GuestID
}

type SearchBookingsRequest struct {
	GetFirst *bool                 `json:"getFirst,omitempty"`
	Select   *[]string             `json:"select,omitempty"`
	Filters  SearchBookingsFilters `json:"filters"`
	Order    *OrderField           `json:"order,omitempty"`
}

type SearchBookingsFilters struct {
	InResa       string      `json:"inResa,omitempty""`
	Booked       bool        `json:"booked,omitempty"`
	IsMasterResa bool        `json:"isMasterResa,omitempty"`
	GuestId      []uuid.UUID `json:"guestId,omitempty"`
	NotGuestId   []uuid.UUID `json:"notGuestId,omitempty"`
}

type OrderField struct {
	ASC bool `json:"asc"`
}

type SearchBookingsResponse struct {
	Bookings []*Booking `json:"bookings"`
}
