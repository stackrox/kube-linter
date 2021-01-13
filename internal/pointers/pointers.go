package pointers

// Bool returns a pointer to a bool.
func Bool(b bool) *bool {
	return &b
}

// Int32 returns a pointer to an int32.
func Int32(i int32) *int32 {
	return &i
}

// Int64 returns a pointer to an int64.
func Int64(i int64) *int64 {
	return &i
}

// Int returns a pointer to an int.
func Int(i int) *int {
	return &i
}
