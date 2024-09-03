package slicesutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, []string{"1", "2", "3"}, Map([]int{1, 2, 3}, func(i int) string {
		return fmt.Sprintf("%d", i)
	}))
}
