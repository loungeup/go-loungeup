package xmltest

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertXMLBody asserts that the given body is equal to the given XML string.
func AssertXMLBody(t *testing.T, body io.ReadCloser, want string) {
	defer body.Close()

	bodyJustRead, err := io.ReadAll(body)
	assert.NoError(t, err)
	assert.Equal(t, inlineString(want), inlineString(string(bodyJustRead)))
}

func inlineString(value string) string {
	result := strings.Clone(value)
	result = strings.ReplaceAll(result, "\n", "")
	result = strings.ReplaceAll(result, "\t", "")

	return result
}
