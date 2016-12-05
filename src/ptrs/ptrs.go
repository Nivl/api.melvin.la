package ptrs

// NewBool takes a bool and returns a pointer to a new bool of the same value
func NewBool(v bool) *bool {
	b := v
	return &b
}
