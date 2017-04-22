package rtb

import (
	"octopus/exchange"
)

type sortedAd []exchange.Advertise

func (sa sortedAd) Len() int {
	return len(sa)
}

func (sa sortedAd) Less(i, j int) bool {
	cpi := sa[i].MaxCPM() * sa[i].Demand().Handicap()
	cpj := sa[j].MaxCPM() * sa[j].Demand().Handicap()
	return cpi > cpj
}

func (sa sortedAd) Swap(i, j int) {
	sa[i], sa[j] = sa[j], sa[i]
}
