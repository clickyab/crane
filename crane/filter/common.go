package filter

import "clickyab.com/exchange/crane/entity"

func hasInt64(slice []int64, elem int64) bool {
	for i := range slice {
		if slice[i] == elem {
			return true
		}
	}
	return false
}

func hasInt(slice []int, elem int) bool {
	if len(slice) == 0 {
		return false
	}
	for i := range slice {
		if slice[i] == elem {
			return true
		}
	}
	return false
}

// hasCategory check for atleast one category to match. one is ok.
func hasCategory(impCat []entity.Category, adCat []entity.Category) bool {
	for _, i := range adCat {
		for _, j := range impCat {
			if i == j {
				return true
			}
		}
	}
	return false
}
