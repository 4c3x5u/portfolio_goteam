package assert

// Comparable is a generic type constraint for all comparable data types in Go.
type Comparable interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 |
		~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

// EqualArr compares two arrays and returns a boolean based on whether or not
// they contain exactly the same items.
func EqualArr[T Comparable](a []T, b []T) bool {
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
