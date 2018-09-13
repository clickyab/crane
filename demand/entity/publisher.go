package entity

import (
	"clickyab.com/crane/openrtb"
)

// PublisherType is the type of publisher
type PublisherType string

const (
	// PublisherTypeApp is the app
	PublisherTypeApp PublisherType = "app"
	// PublisherTypeWeb is the web
	PublisherTypeWeb PublisherType = "web"
)

func (s PublisherType) String() string {
	return string(s)
}

// PublisherAttributes is the key for publisher attributes
type PublisherAttributes int

const (
	// PAMobileAd determine if banner ad show be shown in mobile version or not
	PAMobileAd PublisherAttributes = iota
	// PAFatFinger determine sensitivity of touch ads in app
	PAFatFinger
)

// Publisher is the publisher interface
type Publisher interface {
	ID() int64
	// FloorCPM is the floor cpm for publisher
	FloorCPM() int64
	// Name of publisher
	Name() string
	// MinBid is the minimum CPC requested for this requests
	MinBid() int64
	// Supplier return the exchange object for this publisher
	Supplier() Supplier
	// CTR returns ctr of a slot with specific size
	CTR(int32) float32
	// Type return type of this publisher
	Type() PublisherType
	// Attributes si any other attributes that is not generally required for other part of the system
	Attributes() map[PublisherAttributes]interface{}
	// MinCPC return min cpc based on ad type
	MinCPC(string) float64
	// Categories return categories
	Categories() []openrtb.ContentCategory
}
