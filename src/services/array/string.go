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
