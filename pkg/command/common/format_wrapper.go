package common

import (
	"fmt"

	"golang.stackrox.io/kube-linter/internal/flagutil"
)

// This block enumerates all supported formats in the checks and templates commands.
var (
	PlainFormat    = newFormat("plain")
	MarkdownFormat = newFormat("markdown")
)
var (
	// AllSupportedFormats contains the list of all supported formats.
	AllSupportedFormats []string

	// FormatValueFactory can generate format flag values.
	FormatValueFactory = flagutil.NewEnumValueFactory("output format", AllSupportedFormats)
)

func newFormat(f string) string {
	for _, existingFormat := range AllSupportedFormats {
		if existingFormat == f {
			panic(fmt.Sprintf("duplicate format: %s", f))
		}
	}
	AllSupportedFormats = append(AllSupportedFormats, f)
	return f
}
