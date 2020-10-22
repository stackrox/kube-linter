package utils

// IgnoreError is useful when you want to defer a func that returns an error,
// but ignore the error without having the linter complain.
func IgnoreError(f func() error) {
	_ = f()
}
