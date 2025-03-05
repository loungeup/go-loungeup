package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMessage(t *testing.T) {
	tests := map[string]struct {
		in   error
		want string
	}{
		"no error":      {in: nil, want: ""},
		"unknown error": {in: assert.AnError, want: errorMessageInternal},
		"with message": {
			in: &Error{
				Message: "Could not found the requested resource",
				Code:    CodeNotFound,
				UnderlyingError: &Error{
					Message: "Record not found",
				},
			},
			want: "Could not found the requested resource",
		},
		"with code": {
			in: &Error{
				Code: CodeNotFound,
				UnderlyingError: &Error{
					Message: "Record not found",
				},
			},
			want: errorMessageNotFound,
		},
		"with underlying error": {
			in: &Error{
				UnderlyingError: &Error{
					Message: "Record not found",
				},
			},
			want: "Record not found",
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, ErrorMessage(tt.in))
		})
	}
}
