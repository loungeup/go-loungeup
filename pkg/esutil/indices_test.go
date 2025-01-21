package esutil

import (
	"testing"
	"time"

	"github.com/loungeup/go-loungeup"
	"github.com/stretchr/testify/assert"
)

func TestIndicesMakerDevelopmentAt(t *testing.T) {
	tests := map[string]struct {
		in   time.Time
		want *Indices
	}{
		"before 2020": {
			in: time.Date(2019, 12, 31, 23, 59, 59, 0, time.UTC),
			want: &Indices{
				Bookings: "development-guestbookings-global",
				Guests:   "development-guestcards-global",
			},
		},
		"in 2021": {
			in: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &Indices{
				Bookings: "development-guestbookings-global",
				Guests:   "development-guestcards-global",
			},
		},
		"in 2023": {
			in: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &Indices{
				Bookings: "development-guestbookings-global",
				Guests:   "development-guestcards-global",
			},
		},
		"after 2100": {
			in: time.Date(2101, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &Indices{
				Bookings: "development-guestbookings-global",
				Guests:   "development-guestcards-global",
			},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got := MakeIndices(loungeup.PlatformDevelopment).At(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIndicesMakerStudioAt(t *testing.T) {
	tests := map[string]struct {
		in   time.Time
		want *Indices
	}{
		"before 2020": {
			in: time.Date(2019, 12, 31, 23, 59, 59, 0, time.UTC),
			want: &Indices{
				Bookings: "studio-guestbookings-global",
				Guests:   "studio-guestcards-global",
			},
		},
		"in 2021": {
			in: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &Indices{
				Bookings: "studio-guestbookings-global",
				Guests:   "studio-guestcards-global",
			},
		},
		"in 2023": {
			in: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &Indices{
				Bookings: "studio-guestbookings-global",
				Guests:   "studio-guestcards-global",
			},
		},
		"after 2100": {
			in: time.Date(2101, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &Indices{
				Bookings: "studio-guestbookings-global",
				Guests:   "studio-guestcards-global",
			},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got := MakeIndices(loungeup.PlatformStudio).At(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIndicesMakerProductionAt(t *testing.T) {
	tests := map[string]struct {
		in   time.Time
		want *Indices
	}{
		"before 2020": {
			in: time.Date(2019, 12, 31, 23, 59, 59, 0, time.UTC),
			want: &Indices{
				Bookings: "production-guestbookings-global",
				Guests:   "production-guestcards-global",
			},
		},
		"in 2021": {
			in: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &Indices{
				Bookings: "production-guestbookings-2021",
				Guests:   "production-guestcards-global",
			},
		},
		"in 2023": {
			in: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &Indices{
				Bookings: "production-guestbookings-2023-01",
				Guests:   "production-guestcards-2023-01",
			},
		},
		"after 2100": {
			in: time.Date(2101, 1, 1, 0, 0, 0, 0, time.UTC),
			want: &Indices{
				Bookings: "production-guestbookings-global",
				Guests:   "production-guestcards-2101-01",
			},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got := MakeIndices(loungeup.PlatformProduction).At(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
