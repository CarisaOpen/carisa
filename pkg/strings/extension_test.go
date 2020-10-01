package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings_Concat(t *testing.T) {
	c := Concat("a", "b")
	assert.Equal(t, c, "ab")
}
