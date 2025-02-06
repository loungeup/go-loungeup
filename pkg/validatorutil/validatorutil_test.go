package validatorutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultValidator(t *testing.T) {
	require.NotNil(t, Default())
}
