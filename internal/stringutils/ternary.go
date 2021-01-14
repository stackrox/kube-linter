package stringutils

// Ternary does a ternary based on the condition.
func Ternary(condition bool, ifTrue, ifFalse string) string {
	if condition {
		return ifTrue
	}
	return ifFalse
}
