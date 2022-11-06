package assert

// EqualArr compares two arrays and returns a boolean based on whether or not
// they contain exactly the same items.
func EqualArr[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
