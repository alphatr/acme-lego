package common

// DefaultString 返回默认字符串
func DefaultString(input string, defaultValue string) string {
	if len(input) > 0 {
		return input
	}

	return defaultValue
}
