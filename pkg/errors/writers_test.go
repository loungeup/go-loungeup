package errors

import (
	"io"
	"log/slog"
	"testing"

	"github.com/jirenius/go-res"
	"github.com/stretchr/testify/assert"
)

func TestLogAndWriteRESError(t *testing.T) {
	tests := map[string]struct {
		in          error
		assertLog   func(t *testing.T, msg string, args ...slog.Attr)
		assertWrite func(t *testing.T, err error)
	}{
		"no error": {
			in: nil,
			assertLog: func(t *testing.T, msg string, args ...slog.Attr) {
				assert.Equal(t, "", msg)
				assert.Empty(t, args)
			},
			assertWrite: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		"unknown error": {
			in: io.EOF,
			assertLog: func(t *testing.T, msg string, args ...slog.Attr) {
				assert.Equal(t, io.EOF.Error(), msg)
				assert.Len(t, args, 1)
			},
			assertWrite: func(t *testing.T, err error) {
				assert.EqualError(t, err, errorMessageInternal)
			},
		},
		"LoungeUp error": {
			in: &Error{Code: CodeNotFound, Message: "Could not find the resource", UnderlyingError: io.EOF},
			assertLog: func(t *testing.T, msg string, args ...slog.Attr) {
				assert.Equal(t, io.EOF.Error(), msg)
				assert.Len(t, args, 2)
			},
			assertWrite: func(t *testing.T, err error) {
				assert.EqualError(t, err, "Could not find the resource")
			},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			logger, writer := &errorLoggerMock{}, &errorWriterMock{}

			LogAndWriteRESError(logger, writer, tt.in)

			tt.assertLog(t, logger.msg, logger.args...)
			tt.assertWrite(t, writer.err)
		})
	}
}

func TestGetRESErrorCode(t *testing.T) {
	tests := map[string]struct {
		in   error
		want string
	}{
		"no error":      {in: nil, want: ""},
		"conflict":      {in: &Error{Code: CodeConflict}, want: res.CodeInvalidParams},
		"internal":      {in: &Error{Code: CodeInternal}, want: res.CodeInternalError},
		"invalid":       {in: &Error{Code: CodeInvalid}, want: res.CodeInvalidParams},
		"not found":     {in: &Error{Code: CodeNotFound}, want: res.CodeNotFound},
		"custom error":  {in: &Error{Code: "guestProfile.guestNotFound"}, want: "guestProfile.guestNotFound"},
		"unknown error": {in: io.EOF, want: res.CodeInternalError},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, getRESErrorCode(tt.in))
		})
	}
}

type errorLoggerMock struct {
	msg  string
	args []slog.Attr
}

func (m *errorLoggerMock) Error(msg string, args ...slog.Attr)          { m.msg, m.args = msg, args }
func (m *errorLoggerMock) FormattedError(msg string, args ...slog.Attr) { m.msg, m.args = msg, args }

type errorWriterMock struct {
	err error
}

func (m *errorWriterMock) Error(err error) { m.err = err }
