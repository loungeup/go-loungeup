// This package provides a structured error type and functions to work with it. It is largely based on the work of
// Ben Johnson. Reference: https://middlemost.com/failure-is-your-domain/failure-is-your-domain.pdf.
package errors

import "errors"

// Pre-defined error codes.
const (
	CodeConflict = "conflict"
	CodeInternal = "internal"
	CodeInvalid  = "invalid"
	CodeNotFound = "notFound"
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

	if err, ok := err.(*Error); ok && err.Message != "" {
		return err.Message
	} else if ok && err.UnderlyingError != nil {
		return ErrorMessage(err.UnderlyingError)
	}

	return "An internal error has occurred. Please contact technical support."
}

func Is(err, target error) bool { return errors.Is(err, target) }
