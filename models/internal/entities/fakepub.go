package entities

import (
	"clickyab.com/crane/demand/entity"
	openrtb "clickyab.com/crane/openrtb/v2.5"
)

type fakepub struct {
	name  string
	s     entity.Supplier
	ptype entity.PublisherType

	att map[entity.PublisherAttributes]interface{}
}

// MaxCPC return max allowed cpc for publisher
func (fp *fakepub) MaxCPC() float64 {
	return 0
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

func (fp *fakepub) CTR(int32) float32 {
	return -1
}

// NewFakePublisher return a publisher with fake data
func NewFakePublisher(s entity.Supplier, name string, t entity.PublisherType) entity.Publisher {
	return &fakepub{name: name, s: s, ptype: t}
}
