package timeutil

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatIfNotZero(t *testing.T) {
	tests := map[string]struct {
		in   time.Time
		want string
	}{
		"zero": {
			in:   time.Time{},
			want: "",
		},
		"not zero": {
			in:   time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: "2020-01-01",
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatIfNotZero(time.DateOnly, tt.in))
		})
	}
}

func TestMostRecentAndOldest(t *testing.T) {
	tests := map[string]struct {
		in                         []time.Time
		wantMostRecent, wantOldest time.Time
	}{
		"empty": {
			in:             []time.Time{},
			wantMostRecent: time.Time{},
			wantOldest:     time.Time{},
		},
		"one time": {
			in: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			wantMostRecent: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			wantOldest:     time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		"many times": {
			in: []time.Time{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
			},
			wantMostRecent: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
			wantOldest:     time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.wantMostRecent, MostRecent(tt.in...))
			assert.Equal(t, tt.wantOldest, Oldest(tt.in...))
		})
	}
}

func TestRFC3339Time(t *testing.T) {
	decoded := RFC3339Time{}
	assert.NoError(t, json.Unmarshal(json.RawMessage(`"2025-01-01T00:00:00Z"`), &decoded))

	assert.Equal(t, NewRFC3339Time(time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)), decoded)

	encoded, err := json.Marshal(decoded)
	assert.NoError(t, err)
	assert.Equal(t, `"2025-01-01T00:00:00Z"`, string(encoded))
}

func TestDateOnlyTime(t *testing.T) {
	decoded := DateOnlyTime{}
	assert.NoError(t, json.Unmarshal(json.RawMessage(`"2025-01-01"`), &decoded))

	assert.Equal(t, NewDateOnlyTime(time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)), decoded)

	encoded, err := json.Marshal(decoded)
	assert.NoError(t, err)
	assert.Equal(t, `"2025-01-01"`, string(encoded))
}
