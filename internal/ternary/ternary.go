package ternary

// String does a ternary statement on strings.
func String(condition bool, ifTrue, ifFalse string) string {
	if condition {
		return ifTrue
	}
	return ifFalse
}
