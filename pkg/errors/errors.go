// This package provides a structured error type and functions to work with it. It is largely based on the work of
// Ben Johnson. Reference: https://www.gobeyond.dev/failure-is-your-domain/.
package errors

import "errors"

// Pre-defined error codes.
const (
	CodeConflict = "conflict"
	CodeInternal = "internal"
	CodeInvalid  = "invalid"
	CodeNotFound = "notFound"

	errorMessageConflict = "Conflict"
	errorMessageInternal = "An internal error has occurred. Please contact technical support."
	errorMessageInvalid  = "Invalid"
	errorMessageNotFound = "Not found"
)

type Error struct {
	// Code is machine-readable.
	Code string

	// Message is human-readable.
	Message string

	// Operation that caused the error.
	Operation string

	// UnderlyingError that caused this error, if any.
	UnderlyingError error
}

func (e *Error) Error() string {
	result := ""

	if e.Operation != "" {
		result += e.Operation + ": "
	}

	if e.UnderlyingError != nil {
		result += e.UnderlyingError.Error() + ""
	} else {
		if e.Code != "" {
			result += "<" + e.Code + "> "
		}

		result += e.Message
	}

	return result
}

func (e *Error) defaultMessage() string {
	switch e.Code {
	case CodeConflict:
		return errorMessageConflict
	case CodeInvalid:
		return errorMessageInvalid
	case CodeNotFound:
		return errorMessageNotFound
	default:
		return errorMessageInternal
	}
}

func ErrorCode(err error) string {
	if err == nil {
		return ""
	}

	if err, ok := err.(*Error); ok && err.Code != "" {
		return err.Code
	} else if ok && err.UnderlyingError != nil {
		return ErrorCode(err.UnderlyingError)
	}

	return CodeInternal
}

func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	switch err, ok := err.(*Error); {
	case ok && err.Message != "":
		return err.Message
	case ok && err.Code != "":
		return err.defaultMessage()
	case ok && err.UnderlyingError != nil:
		return ErrorMessage(err.UnderlyingError)
	default:
		return errorMessageInternal
	}
}

func As(err error, target any) bool { return errors.As(err, target) }
func Is(err, target error) bool     { return errors.Is(err, target) }
