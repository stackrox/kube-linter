package flagutil

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	"github.com/pkg/errors"
	"golang.stackrox.io/kube-linter/internal/set"
	"golang.stackrox.io/kube-linter/internal/utils"
)

// EnumFlag allows setting a list of values.
type EnumFlag struct {
	flagDescription string
	allowedValues   set.FrozenStringSet

	currentValue string
}

// String implements pflag.Value.
// It is the preferred method of retrieving the set value for the enum.
func (e *EnumFlag) String() string {
	return e.currentValue
}

// Set implements pflag.Value.
func (e *EnumFlag) Set(input string) error {
	if !e.allowedValues.Contains(input) {
		return errors.Errorf("%q is not a valid option (valid options are: %v)", input, e.getAllowedValuesString())
	}
	e.currentValue = input
	return nil
}

// Type implements pflag.Value.
func (e *EnumFlag) Type() string {
	return "string"
}

// Check that EnumFlag implements pflag.Value interface.
var _ pflag.Value = (*EnumFlag)(nil)

// Usage returns a string that can be used as help text for this flag.
// It will include the flag type and the list of allowed values.
func (e *EnumFlag) Usage() string {
	return fmt.Sprintf("%s. Allowed values: %v.", e.flagDescription, e.getAllowedValuesString())
}

func (e *EnumFlag) getAllowedValuesString() string {
	return strings.Join(e.allowedValues.AsSortedSlice(func(i, j string) bool {
		return i < j
	}), ", ")
}

// NewEnumFlag creates and returns an enum flag value with the given description, allowedValues and defaultValue.
func NewEnumFlag(flagDescription string, allowedValues []string, defaultValue string) *EnumFlag {
	allowedValuesSet := set.NewFrozenStringSet(allowedValues...)
	partialValue := &EnumFlag{flagDescription: flagDescription, allowedValues: allowedValuesSet}
	utils.Must(partialValue.Set(defaultValue))
	return partialValue
}
