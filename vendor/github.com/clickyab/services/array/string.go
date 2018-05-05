package array

// StringInArray return true if array contain the string
func StringInArray(in string, arr ...string) bool {
	for i := range arr {
		if arr[i] == in {
			return true
		}
	}

	return false
}

// IntInArray return true if array contain the number
func IntInArray(in int, arr ...int) bool {
	for i := range arr {
		if arr[i] == in {
			return true
		}
	}

	return false
}
