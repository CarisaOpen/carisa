package strings

import (
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

func Lpad(s1 string, length int, s2 string) string {
	s1l := len(s1)
	if length <= s1l {
		return s1
	}
	return Concat(strings.Repeat(s2, length-s1l), s1)
}
