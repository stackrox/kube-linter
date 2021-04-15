package stringutils

// OrDefault returns the string if it's not empty, or the default.
func OrDefault(s, defaultValue string) string {
	if s != "" {
		return s
	}
	return defaultValue
}

// PointerOrDefault returns the string if it's not nil nor empty, or the default.
func PointerOrDefault(s *string, defaultValue string) string {
	if s == nil {
		return defaultValue
	}

	return OrDefault(*s, defaultValue)
}
