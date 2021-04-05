package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings_Concat(t *testing.T) {
	c := Concat("a", "b")
	assert.Equal(t, c, "ab")
}

func TestStrings_Convert(t *testing.T) {
	c := ConvertBytes([]byte("string"))
	assert.Equal(t, c, "string")
}

func TestStrings_ConvertFromuint64(t *testing.T) {
	c := Convertuint32(uint32(1324432324))
	assert.Equal(t, c, "1324432324")
}

func TestStrings_Lpad(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		len    int
		repeat string
		res    string
	}{
		{
			name:   "lpad - bb",
			input:  "bb",
			len:    1,
			repeat: "a",
			res:    "bb",
		},
		{
			name:   "lpad - bb",
			input:  "bb",
			len:    2,
			repeat: "a",
			res:    "bb",
		},
		{
			name:   "lpad - abb",
			input:  "bb",
			len:    3,
			repeat: "a",
			res:    "abb",
		},
		{
			name:   "lpad - aabb",
			input:  "bb",
			len:    4,
			repeat: "a",
			res:    "aabb",
		},
	}

	for _, tt := range tests {
		res := Lpad(tt.input, tt.len, tt.repeat)
		assert.Equal(t, tt.res, res)
	}
}
