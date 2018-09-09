package entities

import (
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/openrtb"
)

type fakepub struct {
	name  string
	s     entity.Supplier
	ptype entity.PublisherType

	att map[entity.PublisherAttributes]interface{}
}

// Categories return publisher categories
func (fp *fakepub) Categories() []openrtb.ContentCategory {
	panic("implement me")
}

func (fp *fakepub) MinCPC(string) float64 {
	panic("implement me")
}

func (fp *fakepub) Attributes() map[entity.PublisherAttributes]interface{} {
	if fp.att == nil {
		fp.att = make(map[entity.PublisherAttributes]interface{})
	}

	return fp.att
}

func (fp fakepub) Type() entity.PublisherType {
	return fp.ptype
}

func (fakepub) ID() int64 {
	return 0
}

// TODO : Floor and soft floor based on supplier
func (fp *fakepub) FloorCPM() int64 {
	return 0
}

func (fp *fakepub) Name() string {
	return fp.name
}

func (fp *fakepub) MinBid() int64 {
	return fp.s.DefaultMinBid()
}

func (fp *fakepub) Supplier() entity.Supplier {
	return fp.s
}

func (fp *fakepub) CTR(int) float64 {
	return -1
}

// NewFakePublisher return a publisher with fake data
func NewFakePublisher(s entity.Supplier, name string, t entity.PublisherType) entity.Publisher {
	return &fakepub{name: name, s: s, ptype: t}
}
