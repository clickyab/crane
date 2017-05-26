package mocks

import "clickyab.com/exchange/octopus/exchange"

type Supplier struct {
	SName            string
	SFloorCPM        int64
	SSoftFloorCPM    int64
	SExcludedDemands []string
	SShare           int
}

func (s Supplier) Name() string {
	return s.SName
}

func (s Supplier) FloorCPM() int64 {
	return s.SFloorCPM
}

func (s Supplier) SoftFloorCPM() int64 {
	return s.SSoftFloorCPM
}

func (s Supplier) ExcludedDemands() []string {
	return s.SExcludedDemands
}

func (s Supplier) Share() int {
	return s.SShare
}

func (s Supplier) Renderer() exchange.Renderer {
	panic("not needed")
}
