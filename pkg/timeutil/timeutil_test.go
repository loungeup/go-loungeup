package timeutil

import (
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
