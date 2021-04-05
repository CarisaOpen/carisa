package strings

import (
	"strconv"
	"strings"
)

// Concat joins many strings
func Concat(params ...string) string {
	var b strings.Builder
	b.Grow(len(params) * 5)
	for _, str := range params {
		b.WriteString(str)
	}
	return b.String()
}

// ConvertBytes converts array of byte to string without escaping to heap
func ConvertBytes(s []byte) string {
	var b strings.Builder
	b.Grow(len(s))
	b.Write(s)
	return b.String()
}

// Convertuint32 converts a uint64 to string
func Convertuint32(u uint32) string {
	return strconv.FormatUint(uint64(u), 10)
}

func Lpad(s1 string, length int, s2 string) string {
	s1l := len(s1)
	if length <= s1l {
		return s1
	}
	return Concat(strings.Repeat(s2, length-s1l), s1)
}
