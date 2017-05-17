package rtb

import (
	"octopus/exchange"
	"services/dset"
	"sort"
	"time"
)

// SelectCPM is the simplest way to bid. sort the value, return the
func SelectCPM(imp exchange.Impression, all map[string][]exchange.Advertise) (res map[string]exchange.Advertise) {
	res = make(map[string]exchange.Advertise, len(all))

	set := dset.NewDistributedSet("EXC" + imp.Source().Supplier().Name() + imp.PageTrackID())
	for id := range all {
		this := moderate(imp.Source(), all[id])
		if len(this) == 0 {
			res[id] = nil
			continue
		}

		sorted := sortedAd(rmDuplicate(set, this))
		sort.Sort(sorted)

		res[id] = sorted[0]
		lower := imp.Source().SoftFloorCPM()
		if lower > res[id].MaxCPM() {
			lower = imp.Source().FloorCPM()
		}
		if len(sorted) > 1 && sorted[1].MaxCPM() > lower {
			lower = sorted[1].MaxCPM()
		}

		if lower < res[id].MaxCPM() {
			res[id].SetWinnerCPM(lower + 1)
		} else {
			res[id].SetWinnerCPM(res[id].MaxCPM())
		}

		set.Add(res[id].ID())
	}

	set.Save(time.Minute)
	return res
}

// moderate remove unacceptable ads for publisher
func moderate(imp exchange.Rater, ads []exchange.Advertise) []exchange.Advertise {
	rds := make([]exchange.Advertise, 0)
	for _, ad := range ads {
		if reduce(imp, ad) {
			rds = append(rds, ad)
		}
	}

	return rds
}

func rmDuplicate(set dset.DistributedSet, ads []exchange.Advertise) []exchange.Advertise {
	all := set.Members()
	var res []exchange.Advertise
bigLoop:
	for id := range ads {
		for _, adID := range all {
			if adID == ads[id].ID() {
				continue bigLoop
			}
		}
		res = append(res, ads[id])
	}
	return res
}

func reduce(pub exchange.Rater, ad exchange.Rater) bool {
	p := pub.Rates()
	a := ad.Rates()
	for _, ar := range a {
		for _, pr := range p {
			if ar == pr {
				return false
			}
		}
	}
	return true
}
