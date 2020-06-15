package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcat(t *testing.T) {
	c := Concat("a", "b")
	assert.Equal(t, c, "ab")
}
