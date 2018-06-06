package entity

import (
	"database/sql/driver"
	"fmt"

	"github.com/clickyab/services/array"
)

// PublisherType is the type of publisher
type PublisherType string

const (
	// PublisherTypeApp is the app
	PublisherTypeApp PublisherType = "app"
	// PublisherTypeWeb is the web
	PublisherTypeWeb PublisherType = "web"
)

func (e PublisherType) String() string {
	return string(e)
}

// Scan convert the json array ino string slice
func (e *PublisherType) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return fmt.Errorf("unsupported type")
	}
	if !PublisherType(b).IsValid() {
		return fmt.Errorf("invalid value")
	}
	*e = PublisherType(b)
	return nil
}

// Value try to get the string slice representation in database
func (e PublisherType) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, fmt.Errorf("invalid publisher type: %s", e)
	}

	return string(e), nil
}

// IsValid try to validate enum value on ths type
func (e PublisherType) IsValid() bool {
	return array.StringInArray(
		string(e),
		string(PublisherTypeApp),
		string(PublisherTypeWeb),
	)
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
	CTR(int) float64
	// Type return type of this publisher
	Type() PublisherType
	// Attributes si any other attributes that is not generally required for other part of the system
	Attributes() map[PublisherAttributes]interface{}
	// MinCPC return min cpc based on ad type
	MinCPC(string) float64
	// Categories return categories
	Categories() []string
}
