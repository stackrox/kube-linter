package utils

// Must panics if any of the errors are not nil.
// It is intended for use in cases where an error returned would
// mean a programming error.
func Must(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
