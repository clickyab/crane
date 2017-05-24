package filter

import "clickyab.com/exchange/crane/entity"

// Target if the campaign target and imp target is ok
func Target(imp entity.Impression, ad entity.Advertise) bool {
	t := imp.Source().AcceptedTarget()
	var r int
	switch t {
	case entity.TargetApp:
		r = int(entity.TargetApp)
	case entity.TargetVast:
		r = int(entity.TargetVast)
	case entity.TargetWeb:
		r = int(entity.TargetWeb)
	}

	for _, i := range ad.Campaign().Target() {
		if int(i) == r {
			return true
		}
	}

	return false
}
