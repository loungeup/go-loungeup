package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBookingInStayDates(t *testing.T) {
	tests := map[string]struct {
		in   *Booking
		want []time.Time
	}{
		"simple": {
			in: &Booking{
				Arrival:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				Departure: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
			},
			want: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		"only arrival": {
			in: &Booking{
				Arrival: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: []time.Time{
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.in.InStayDates())
		})
	}
}
