package jetstreamutil

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThrottler(t *testing.T) {
	msg := &msgMock{
		data: []byte(`{"params": {"throttlerInterval": "5ms"}}`),
	}

	NewThrottler(
		WithThrottlerInterval(time.Minute), // Should be ignored.
		WithThrottlerInProgressInterval(time.Millisecond),
	).Handle(func(msg jetstream.Msg) { msg.Ack() })(msg)

	time.Sleep(6 * time.Millisecond)

	require.Equal(t, 1, msg.ackCount)
	require.Equal(t, 4, msg.inProgressCount)
}

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

type msgMock struct {
	subject string
	data    []byte

	ackCount        int
	inProgressCount int
}

var _ (jetstream.Msg) = (*msgMock)(nil)

func (m *msgMock) Ack() error                                { m.ackCount++; return nil }
func (m *msgMock) Data() []byte                              { return m.data }
func (m *msgMock) DoubleAck(context.Context) error           { return nil }
func (m *msgMock) Headers() nats.Header                      { return nil }
func (m *msgMock) InProgress() error                         { m.inProgressCount++; return nil }
func (m *msgMock) Metadata() (*jetstream.MsgMetadata, error) { return nil, nil }
func (m *msgMock) Nak() error                                { return nil }
func (m *msgMock) NakWithDelay(delay time.Duration) error    { return nil }
func (m *msgMock) Reply() string                             { return "" }
func (m *msgMock) Subject() string                           { return m.subject }
func (m *msgMock) Term() error                               { return nil }
func (m *msgMock) TermWithReason(reason string) error        { return nil }
