package resutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeferredCallRequest(t *testing.T) {
	request := &DeferredCallRequest{}

	request.Error(assert.AnError)
	assert.Equal(t, assert.AnError, request.GetError())

	request.Resource("authority.countries.fr")
	assert.Equal(t, "authority.countries.fr", request.GetRID())

	request.OK("foo")
	assert.Equal(t, "foo", request.GetResult())
}
