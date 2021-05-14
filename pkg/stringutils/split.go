package stringutils

import (
	"strings"
)

// Split2 splits the given string at the given separator, returning the part before and after the separator as two
// separate return values.
// If the string does not contain `sep`, the entire string is returned as the first return value.
func Split2(str, sep string) (string, string) {
	splitIdx := strings.Index(str, sep)
	if splitIdx == -1 {
		return str, ""
	}
	return str[:splitIdx], str[splitIdx+len(sep):]
}
