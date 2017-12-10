package filter

import "clickyab.com/crane/crane/entity"

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

func hasInt(empty bool, slice []int, elem int) bool {
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

// hasCategory check for atleast one category to match. one is ok.
func hasCategory(empty bool, impCat []entity.Category, adCat []entity.Category) bool {
	if len(impCat) == 0 || len(adCat) == 0 {
		return empty
	}
	for _, i := range adCat {
		for _, j := range impCat {
			if i == j {
				return true
			}
		}
	}
	return false
}
