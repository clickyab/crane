package filter

func hasString(empty bool, slice []string, elem string) bool {
	if len(slice) == 0 {
		return empty
	}
	for i := range slice {
		if slice[i] == elem {
			return true
		}
	}
	return false
}
