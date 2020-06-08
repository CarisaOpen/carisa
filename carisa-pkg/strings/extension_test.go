package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConcat(t *testing.T) {
	c := Concat("a", "b")
	assert.Equal(t, c, "ab")
}
