package stringutils

import (
	"strings"
)

// Repeat repeats the given string `n` times efficiently.
func Repeat(s string, n int) string {
	var sb strings.Builder
	sb.Grow(len([]byte(s)) * n)
	for i := 0; i < n; i++ {
		sb.WriteString(s)
	}
	return sb.String()
}
