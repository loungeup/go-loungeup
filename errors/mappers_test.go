package errors

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestMapSQLError(t *testing.T) {
	tests := map[string]struct{ in, want error }{
		"no error": {in: nil, want: nil},
		"not found": {
			in:   sql.ErrNoRows,
			want: &Error{Code: CodeNotFound, UnderlyingError: sql.ErrNoRows},
		},
		"invalid": {
			in:   &pq.Error{Code: "22000"},
			want: &Error{Code: CodeInvalid, UnderlyingError: &pq.Error{Code: "22000"}},
		},
		"conflict": {
			in:   &pq.Error{Code: "23000"},
			want: &Error{Code: CodeConflict, UnderlyingError: &pq.Error{Code: "23000"}},
		},
		"unknown": {
			in:   fmt.Errorf("unknown error"),
			want: &Error{Code: CodeInternal, UnderlyingError: fmt.Errorf("unknown error")},
		},
		"unknown pq error": {
			in:   &pq.Error{Code: "unknown"},
			want: &Error{Code: CodeInternal, UnderlyingError: &pq.Error{Code: "unknown"}},
		},
	}

	for test, tt := range tests {
		t.Run(test, func(t *testing.T) {
			assert.Equal(t, tt.want, MapSQLError(tt.in))
		})
	}
}
