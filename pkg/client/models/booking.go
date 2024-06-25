package models

import (
	"strconv"
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
	CustomFields       CustomFields `json:"customFields"`
}

func (b *Booking) InStayDates() []time.Time {
	result := []time.Time{b.arrivalDay()}
	for d := b.arrivalDay().Add(aDay); d.Before(b.departureDay()); d = d.Add(aDay) {
		result = append(result, d)
	}

	return result
}

func (b *Booking) arrivalDay() time.Time { return b.Arrival.Truncate(aDay) }

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

func (c CustomFields) Get(key string) string {
	if field, ok := c[key]; ok {
		return field.Value
	}

	return ""
}
