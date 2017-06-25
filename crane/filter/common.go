package filter

import "clickyab.com/crane/crane/entity"

func hasString(slice []string, elem string) bool {
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
