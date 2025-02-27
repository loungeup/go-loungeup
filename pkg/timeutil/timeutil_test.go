package timeutil

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMostRecent(t *testing.T) {
	tests := map[string]struct {
		in   []time.Time
		want time.Time
	}{
		"no times": {in: []time.Time{}, want: time.Time{}},
		"one time": {
			in: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		"many times": {
			in: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, MostRecent(tt.in...))
		})
	}
}

func TestOldest(t *testing.T) {
	tests := map[string]struct {
		in   []time.Time
		want time.Time
	}{
		"no times": {in: []time.Time{}, want: time.Time{}},
		"one time": {
			in: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		"many times": {
			in: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, Oldest(tt.in...))
		})
	}
}

func TestDateUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		in        json.RawMessage
		want      Date
		wantError bool
	}{
		"valid date": {
			in:   json.RawMessage(`"2025-01-24"`),
			want: NewDate(time.Date(2025, 1, 24, 0, 0, 0, 0, time.UTC)),
		},
		"invalid format": {
			in:        json.RawMessage(`"01-01-2020"`),
			wantError: true,
		},
		"invalid date": {
			in:        json.RawMessage(`"2020-13-01"`),
			wantError: true,
		},
		"empty string": {
			in:   json.RawMessage(`""`),
			want: NewDate(time.Time{}),
		},
		"not a string": {
			in:        json.RawMessage(`123`),
			wantError: true,
		},
		"no data": {
			in:   json.RawMessage(`null`),
			want: NewDate(time.Time{}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var got Date

			err := json.Unmarshal(tt.in, &got)
			if tt.wantError {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDateMarshalJSON(t *testing.T) {
	tests := map[string]struct {
		in   Date
		want json.RawMessage
	}{
		"simple": {
			in:   NewDate(time.Date(2025, 1, 24, 12, 30, 0, 0, time.UTC)),
			want: json.RawMessage(`"2025-01-24"`),
		},
		"zero date": {
			in:   NewDate(time.Time{}),
			want: json.RawMessage(`"0001-01-01"`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := json.Marshal(Date(tt.in))
			assert.NoError(t, err)
			assert.Equal(t, tt.want, json.RawMessage(got))
		})
	}
}
