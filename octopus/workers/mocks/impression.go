package mocks

import (
	"net"
	"time"

	"clickyab.com/exchange/octopus/exchange"
)

type Impression struct {
	ITrackID     string
	IIP          net.IP
	ISchema      string
	IUserTrackID string
	IPageTrackID string
	IUserAgent   string
	ISource      Publisher
	ILocation    Location
	IAttributes  map[string]interface{}
	ISlots       []Slot
	ICategory    []exchange.Category
	IPlatform    exchange.ImpressionPlatform
	IUnderFloor  bool
	ITime        time.Time
}

func (i Impression) TrackID() string {
	return i.ITrackID
}

func (i Impression) IP() net.IP {
	return i.IIP
}

func (i Impression) Scheme() string {
	return i.ISchema
}

func (i Impression) UserTrackID() string {
	return i.IUserTrackID
}

func (i Impression) PageTrackID() string {
	return i.IPageTrackID
}

func (i Impression) UserAgent() string {
	return i.IUserAgent
}

func (i Impression) Source() exchange.Publisher {
	return i.ISource
}

func (i Impression) Location() exchange.Location {
	return i.ILocation
}

func (i Impression) Attributes() map[string]interface{} {
	return i.IAttributes
}

func (i Impression) Slots() []exchange.Slot {
	res := make([]exchange.Slot, len(i.ISlots))
	for j := range i.ISlots {
		res[j] = i.ISlots[j]
	}
	return res
}

func (i Impression) Category() []exchange.Category {
	return i.ICategory
}

func (i Impression) Platform() exchange.ImpressionPlatform {
	return i.IPlatform
}

func (i Impression) UnderFloor() bool {
	return i.IUnderFloor
}

func (i Impression) Time() time.Time {
	return i.ITime
}
