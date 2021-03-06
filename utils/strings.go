package utils

// IsStringOfLength returns true if length of the string == n
func IsStringOfLength(msg string, n int) bool {
	if len(msg) != n {
		return false
	}

	return true
}

// IsStringEmpty returns true if the string is empty
func IsStringEmpty(msg string) bool {
	return IsStringOfLength(msg, 0)
}

// ContainsString returns true if the slice contains str.
func ContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}

	return false
}
