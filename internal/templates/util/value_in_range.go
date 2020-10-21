package util

// ValueInRange returns whether the given quantity is in the range between the lowerBound and the upperBound (inclusive).
// A nil upper bound is interpreted as infinity.
func ValueInRange(value, lowerBound int, upperBound *int) bool {
	if value < lowerBound {
		return false
	}
	if upperBound != nil && value > *upperBound {
		return false
	}
	return true
}
