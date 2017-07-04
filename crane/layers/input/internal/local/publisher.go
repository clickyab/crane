package local

import (
	"clickyab.com/crane/crane/entity"
	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/config"
)

const acceptedTarget string = "accepted_target"

var originalUnderFloor = config.RegisterBoolean("crane.input.rest.under_floor", false, "its used when publisher's underfloor isn't set")

// Publisher Publisher
type Publisher struct {
	// Name of publisher
	FName string `json:"name"`
	// SoftFloorCPM is the soft version of floor cpm. if the publisher ahs it, then the system
	// FloorCPM is the floor cpm for publisher
	FFloorCPM int64 `json:"floor_cpm"`
	// try to use this as floor, but if this is not available, the FloorCPM is used
	FSoftFloorCPM int64 `json:"soft_floor_cpm"`
	// Attributes is the generic attribute system
	FAttributes map[string]interface{} `json:"attributes"`
	// Supplier the supplier
	FSupplier string `json:"supplier"`
	// UnderFloor asd
	FUnderFloor *bool
}

// FloorCPM is publisher's cpm floor
func (rp *Publisher) FloorCPM() int64 {
	return rp.FFloorCPM
}

// SoftFloorCPM is publisher soft floor cpm
// @required FSoftFloorCPM
func (rp *Publisher) SoftFloorCPM() int64 {
	return rp.FSoftFloorCPM
}

// UnderFloor is publisher underfloor
func (rp *Publisher) UnderFloor() bool {
	if rp.UnderFloor == nil {
		return originalUnderFloor.Bool()
	}
	return *rp.FUnderFloor
}

// Name is publisher name
// @required FName
func (rp *Publisher) Name() string {
	return rp.FName
}

// AcceptedTargets is publisher's target (web, vast, app, native)
func (rp *Publisher) AcceptedTargets() []entity.Target {
	t, ok := rp.FAttributes[acceptedTarget].([]entity.Target)
	if !ok {
		return []entity.Target{entity.TargetInvalid}
	}

	return t
}

// Attributes is publisher's Attributes
func (rp *Publisher) Attributes() map[string]interface{} {
	return rp.FAttributes
}

// BIDType is publisher's bid type, rest is cpm
func (rp *Publisher) BIDType() entity.BIDType {
	return entity.BIDTypeCPM
}

// MinCPC is publisher's minimum cpc
func (rp *Publisher) MinCPC() int64 {
	logrus.Panic("rest type shouldn't have minCPC")
	return 0
}

// AcceptedTypes is publisher's accepted types (dyn, banner, video, html, native)
func (rp *Publisher) AcceptedTypes() []entity.AdType {
	at, ok := rp.FAttributes["accepted_types"].([]entity.AdType)
	if !ok {
		return nil
	}

	return at
}

// Supplier is publisher's supplier
// @required FSupplier
func (rp *Publisher) Supplier() string {
	return rp.FSupplier
}
