package stringutils

// OrDefault returns the string if it's not empty, or the default.
func OrDefault(s, defaul string) string {
	if s != "" {
		return s
	}
	return defaul
}

// PointerOrDefault returns the string if it's not nil nor empty, or the default.
func PointerOrDefault(s *string, defaul string) string {
	if s == nil {
		return defaul
	}

	return OrDefault(*s, defaul)
}
