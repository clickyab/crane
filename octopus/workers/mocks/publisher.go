package mocks

import (
	"clickyab.com/exchange/octopus/exchange"
)

type Publisher struct {
	PName         string
	PFloorCPM     int64
	PSoftFloorCPM int64
	PAttributes   map[string]interface{}
	PSupplier     Supplier
	PRates        []exchange.Rate
}

func (p Publisher) Name() string {
	return p.PName
}

func (p Publisher) FloorCPM() int64 {
	return p.PFloorCPM
}

func (p Publisher) SoftFloorCPM() int64 {
	return p.PSoftFloorCPM
}

func (p Publisher) Attributes() map[string]interface{} {
	return p.PAttributes
}

func (p Publisher) Supplier() exchange.Supplier {
	return p.PSupplier
}

func (p Publisher) Rates() []exchange.Rate {
	return p.PRates
}
