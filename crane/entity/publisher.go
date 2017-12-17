package entity

const (
	// PublisherTypeApp is the app
	PublisherTypeApp = 1
	// PublisherTypeWeb is the web
	PublisherTypeWeb = 2
	// PublisherTypeVast is the vast
	PublisherTypeVast = 3
)

// BIDType is the bid type for this imp cpc or cpm
type BIDType string

const (
	// BIDTypeCPC is the cost per click type
	BIDTypeCPC BIDType = "CPC"
	//BIDTypeCPM is the cost per view type
	BIDTypeCPM BIDType = "CPM"
)

// Publisher is the publisher interface
type Publisher interface {
	ID() int64
	// PublicID is a fucking shit from the old time, should remove this as soon as possible
	PublicID() string
	// FloorCPM is the floor cpm for publisher
	FloorCPM() int64
	// SoftFloorCPM is the soft version of floor cpm. if the publisher ahs it, then the system
	// try to use this as floor, but if this is not available, the FloorCPM is used
	SoftFloorCPM() int64
	// Name of publisher
	Name() string
	// BIDType return this publisher bid type
	BIDType() BIDType
	// MinBid is the minimum CPC requested for this requests
	MinBid() int64
	// Supplier return the exchange object for this publisher
	Supplier() Supplier
	// CTR returns ctr of a slot with specific size
	CTR(int) float64
}
