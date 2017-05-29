package mocks

import (
	"context"
	"net/http"

	"clickyab.com/exchange/octopus/exchange"
)

type Demand struct {
	DName               string
	DCallRate           int
	DHandicap           int64
	DWhitelistCountries []string
}

func (d *Demand) Name() string {
	return d.DName
}

func (*Demand) Provide(context.Context, exchange.Impression, chan exchange.Advertise) {
	panic("implement me")
}

func (*Demand) Win(context.Context, string, int64) {
	panic("implement me")
}

func (*Demand) Status(context.Context, http.ResponseWriter, *http.Request) {
	panic("implement me")
}

func (d *Demand) Handicap() int64 {
	return d.DHandicap
}

func (d *Demand) CallRate() int {
	return d.DCallRate
}

func (d *Demand) WhiteListCountries() []string {
	return d.DWhitelistCountries
}

func (*Demand) ExcludedSuppliers() []string {
	panic("implement me")
}
