package strings

import "strings"

// Concat joins many strings
func Concat(params ...string) string {
	var b strings.Builder
	b.Grow(len(params) * 5)
	for _, str := range params {
		b.WriteString(str)
	}
	return b.String()
}
