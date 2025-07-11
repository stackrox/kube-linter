package errorhelpers

import (
	"fmt"
	"strings"
)

// ErrorList is a wrapper around many errors
type ErrorList struct {
	start  string
	errors []string
}

// NewErrorList returns a new ErrorList
func NewErrorList(start string) *ErrorList {
	return &ErrorList{
		start: start,
	}
}

// NewErrorListWithErrors returns a new ErrorList with the given errors.
func NewErrorListWithErrors(start string, errors []error) *ErrorList {
	errorList := NewErrorList(start)
	for _, err := range errors {
		errorList.AddError(err)
	}
	return errorList
}

// AddError adds the passed error to the list of errors if it is not nil
func (e *ErrorList) AddError(err error) {
	if err == nil {
		return
	}
	e.errors = append(e.errors, err.Error())
}

// AddErrors adds the non-nil errors in the given slice to the list of errors.
func (e *ErrorList) AddErrors(errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}
		e.errors = append(e.errors, err.Error())
	}
}

// AddWrap is a convenient wrapper around `AddError(fmt.Errorf("%s: %w", msg, err))`.
func (e *ErrorList) AddWrap(err error, msg string) {
	e.AddError(fmt.Errorf("%s: %w", msg, err))
}

// AddWrapf is a convenient wrapper around `AddError(fmt.Errorf(format+": %w", append(args, err)...))`.
func (e *ErrorList) AddWrapf(err error, format string, args ...interface{}) {
	e.AddError(fmt.Errorf(format+": %w", append(args, err)...))
}

// AddString adds a string based error to the list
func (e *ErrorList) AddString(err string) {
	e.errors = append(e.errors, err)
}

// AddStringf adds a templated string
func (e *ErrorList) AddStringf(t string, args ...interface{}) {
	e.errors = append(e.errors, fmt.Sprintf(t, args...))
}

// AddStrings adds multiple string based errors to the list.
func (e *ErrorList) AddStrings(errs ...string) {
	e.errors = append(e.errors, errs...)
}

// ToError returns an error if there were errors added or nil
func (e *ErrorList) ToError() error {
	switch len(e.errors) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("%s error: %s", e.start, e.errors[0])
	default:
		return fmt.Errorf("%s errors: [%s]", e.start, strings.Join(e.errors, ", "))
	}
}

// String converts the list to a string, returning empty if no errors were added.
func (e *ErrorList) String() string {
	err := e.ToError()
	if err == nil {
		return ""
	}
	return err.Error()
}

// ErrorStrings returns all the error strings in this ErrorList as a slice, ignoring the start string.
func (e *ErrorList) ErrorStrings() []string {
	return e.errors
}
