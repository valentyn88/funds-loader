package tools

// IsAmountExceeded compares two amounts.
func IsAmountExceeded(limit, value float64) bool {
	return value > limit
}
