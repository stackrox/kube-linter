package common

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/set"
)

// This block enumerates all supported formats.
var (
	PlainFormat    = newFormat("plain")
	MarkdownFormat = newFormat("markdown")
)
var (
	// AllSupportedFormats contains the list of all supported formats.
	AllSupportedFormats = set.NewStringSet()
)

// FormatWrapper wraps a format flag, implementing the flag.Value interface,
// and enforcing that the value is one of the supported formats.
type FormatWrapper struct {
	Format string
}

// String implements flag.Value.
func (f *FormatWrapper) String() string {
	return f.Format
}

// Set implements flag.Value.
func (f *FormatWrapper) Set(input string) error {
	if !AllSupportedFormats.Contains(input) {
		return errors.Errorf("%q is not a valid option (valid options are %v)", input, AllSupportedFormats.AsSortedSlice(func(i, j string) bool {
			return i < j
		}))
	}
	f.Format = input
	return nil
}

// Type implements flag.Value.
func (f *FormatWrapper) Type() string {
	return "output format"
}

func newFormat(f string) string {
	if !AllSupportedFormats.Add(f) {
		panic(fmt.Sprintf("duplicate format: %s", f))
	}
	return f
}