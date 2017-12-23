package entities

import "clickyab.com/crane/demand/entity"

type fakepub struct {
	name  string
	s     entity.Supplier
	ptype entity.PublisherType
}

func (fp fakepub) Type() entity.PublisherType {
	return fp.ptype
}

func (fakepub) ID() int64 {
	return 0
}

// TODO : Floor and soft floor based on supplier
func (fp *fakepub) FloorCPM() int64 {
	return fp.s.DefaultFloorCPM()
}

func (fp *fakepub) SoftFloorCPM() int64 {
	return fp.s.DefaultSoftFloorCPM()
}

func (fp *fakepub) Name() string {
	return fp.name
}

func (fp *fakepub) BIDType() entity.BIDType {
	return fp.s.BidType()
}

func (fp *fakepub) MinBid() int64 {
	return fp.s.DefaultMinBid()
}

func (fp *fakepub) Supplier() entity.Supplier {
	return fp.s
}

func (fp *fakepub) CTR(int) float64 {
	return fp.s.DefaultCTR()
}

// NewFakePublisher return a publisher with fake data
func NewFakePublisher(s entity.Supplier, name string, t entity.PublisherType) entity.Publisher {
	return &fakepub{name: name, s: s, ptype: t}
}
