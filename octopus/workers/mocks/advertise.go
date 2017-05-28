package mocks

import "clickyab.com/exchange/octopus/exchange"

//advertiser

type Advertiser struct {
	MID          string
	MMaxCPM      int64
	MWidth       int
	MHeight      int
	MURL         string
	MLanding     string
	MSlotTrackID string
	MDemand      exchange.Demand
}

func (a Advertiser) ID() string {
	return a.MID
}

func (a Advertiser) MaxCPM() int64 {
	return a.MMaxCPM
}

func (a Advertiser) Width() int {
	return a.MWidth
}

func (a Advertiser) Height() int {
	return a.MHeight
}

func (a Advertiser) URL() string {
	return a.MURL
}

func (a Advertiser) Landing() string {
	return a.MLanding
}

func (a Advertiser) SlotTrackID() string {
	return a.MSlotTrackID
}

func (Advertiser) Rates() []exchange.Rate {
	panic("implement me")
}

func (Advertiser) TrackID() string {
	panic("implement me")
}

func (Advertiser) SetWinnerCPM(int64) {
	panic("implement me")
}

func (Advertiser) WinnerCPM() int64 {
	panic("implement me")
}

func (a Advertiser) Demand() exchange.Demand {
	return a.MDemand
}
