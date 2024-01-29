package errors

import (
	"database/sql"

	"github.com/lib/pq"
)

// MapSQLError into a well-defined error.
func MapSQLError(err error) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return &Error{
			Code:            CodeNotFound,
			UnderlyingError: err,
		}
	}

	parsedError := &pq.Error{}
	if !As(err, &parsedError) {
		return mapError(err)
	}

	// https://www.postgresql.org/docs/9.3/errcodes-appendix.html
	switch parsedError.Code.Class() {
	case "22":
		return &Error{
			Code:            CodeInvalid,
			UnderlyingError: err,
		}
	case "23":
		return &Error{
			Code:            CodeConflict,
			UnderlyingError: err,
		}
	}

	return mapError(err)
}

// mapError into a well-defined error.
func mapError(err error) error {
	if err == nil {
		return nil
	}

	return &Error{
		Code:            CodeInternal,
		UnderlyingError: err,
	}
}
