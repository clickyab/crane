package mocks

import "clickyab.com/exchange/octopus/exchange"

type Ads struct {
	AWidth       int
	AHeight      int
	AMaxCPM      int64
	ADemand      exchange.Demand
	AID          string
	ALanding     string
	ARates       []exchange.Rate
	ATrackID     string
	AURL         string
	AWinnerCPM   int64
	ASlotTrackID string
}

func (a *Ads) ID() string {
	return a.AID
}

func (a *Ads) MaxCPM() int64 {
	return a.AMaxCPM
}

func (a *Ads) Width() int {
	return a.AWidth
}

func (a *Ads) Height() int {
	return a.AHeight
}

func (a *Ads) URL() string {
	return a.AURL
}

func (a *Ads) Landing() string {
	return a.ALanding
}

func (a *Ads) SlotTrackID() string {
	return a.ASlotTrackID
}

func (a *Ads) Rates() []exchange.Rate {
	return a.ARates
}

func (a *Ads) TrackID() string {
	return a.ATrackID
}

func (*Ads) SetWinnerCPM(int64) {
	panic("implement me")
}

func (a *Ads) WinnerCPM() int64 {
	return a.AWinnerCPM
}

func (a *Ads) Demand() exchange.Demand {
	return a.ADemand
}
