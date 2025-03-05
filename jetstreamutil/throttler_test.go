package jetstreamutil

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThrottlerParamsUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		in   []byte
		want *throttlerParams
	}{
		"simple": {
			in:   []byte(`{"throttlerInterval":"1s"}`),
			want: &throttlerParams{throttlerInterval: time.Second},
		},
		"empty": {
			in:   []byte(`{"throttlerInterval":""}`),
			want: &throttlerParams{throttlerInterval: -1},
		},
		"zero": {
			in:   []byte(`{"throttlerInterval":"0"}`),
			want: &throttlerParams{throttlerInterval: 0},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			got := &throttlerParams{}
			assert.NoError(t, json.Unmarshal(tt.in, got))
			assert.Equal(t, tt.want, got)
		})
	}
}
