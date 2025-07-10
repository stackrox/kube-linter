package errors

import (
	"errors"
	"fmt"
)

func Wrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func Wrapf(err error, format string, args ...any) error {
	return fmt.Errorf(format + ": %w)", append(args, err)...)
}

func Errorf(format string, args ...any) error {
	return fmt.Errorf(format, args)
}

func New(msg string) error {
	return errors.New(msg)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

