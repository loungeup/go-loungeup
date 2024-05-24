package log

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	assert.NotNil(t, Default())
	assert.NotNil(t, Default().Adapter)
	assert.NotPanics(t, func() { Default().Debug("A debug message") })
	assert.NotPanics(t, func() { Default().Error("An error message") })
	assert.NotPanics(t, func() { Default().FormattedDebug("A formatted debug message") })
	assert.NotPanics(t, func() { Default().FormattedError("A formatted error message") })
}

func TestReplaceLogAttribute(t *testing.T) {
	tests := map[string]struct {
		in, want slog.Attr
	}{
		"level": {
			in:   slog.String(slog.LevelKey, "INFO"),
			want: slog.String(statusKey, "info"),
		},
		"message": {
			in:   slog.String(slog.MessageKey, "Test message"),
			want: slog.String(messageKey, "Test message"),
		},
		"time": {
			in:   slog.String(slog.TimeKey, "2022-01-01T12:00:00Z"),
			want: slog.String(timestampKey, "2022-01-01T12:00:00Z"),
		},
		"other": {
			in:   slog.String("error", "An error occurred"),
			want: slog.String("error", "An error occurred"),
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, replaceLogAttribute(nil, tt.in))
		})
	}
}
