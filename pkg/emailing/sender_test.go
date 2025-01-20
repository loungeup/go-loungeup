package emailing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWellKnownEmailSenders(t *testing.T) {
	assert.Equal(t, "support@loungeup.com", WellKnownEmailSenders().LoungeUpSupport)
}
