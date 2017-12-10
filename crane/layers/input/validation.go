package input

import "github.com/bsm/openrtb"

// different impression types validation
func requestValidation(req openrtb.BidRequest) bool {
	for i := range req.Imp {
		if (req.Imp[i].Banner == nil) != (req.Imp[0].Banner == nil) {
			return false
		}
		if (req.Imp[i].Video == nil) != (req.Imp[0].Video == nil) {
			return false
		}
		if (req.Imp[i].Native == nil) != (req.Imp[0].Native == nil) {
			return false
		}
	}
	return true
}
