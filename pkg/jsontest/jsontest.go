package jsontest

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertJSONBody asserts that the given body is equal to the given JSON string.
func AssertJSONBody(t *testing.T, body io.ReadCloser, want string) {
	defer body.Close()

	bodyJustRead, err := io.ReadAll(body)
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(bodyJustRead))
}

// SafeMarshal the given value to JSON, ignoring any errors.
func SafeMarshal(value any) []byte {
	result, _ := json.Marshal(value)

	return result
}

// Stringify returns the JSON-encoded version of the given value as a string.
func Stringify(value any) string {
	return string(SafeMarshal(value))
}
