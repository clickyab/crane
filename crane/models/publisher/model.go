package publisher

import (
	"time"

	"clickyab.com/crane/crane/entity"
)

type (
	// ActiveStatus active field for db
	ActiveStatus string
	// Platforms consist of app, web, vast
	Platforms string
)

const (
	// ActiveStatusTrue ActiveStatusTrue
	ActiveStatusTrue ActiveStatus = "yes"
	// ActiveStatusFalse ActiveStatusFalse
	ActiveStatusFalse ActiveStatus = "no"

	// AppPlatform AppPlatform
	AppPlatform Platforms = "app"
	// VastPlatform VastPlatform
	VastPlatform Platforms = "vast"
	// WebPlatform WebPlatform
	WebPlatform Platforms = "web"
)

// Publisher user model in database
// @Model {
//		table = publishers
//		primary = true, id
//		find_by = id
//		transaction = insert
//		list = yes
// }
type Publisher struct {
	FID           int64        `json:"id" db:"id"`
	UserID        int64        `json:"user_id" db:"user_id"`
	FFloorCPM     int64        `json:"floor_cpm" db:"floor_cpm"`
	FSoftFloorCPM int64        `json:"soft_floor_cpm" db:"soft_floor_cpm"`
	FName         string       `json:"name" db:"name"`
	BidType       int64        `json:"bid_type" db:"bid_type"`
	FUnderFloor   int64        `json:"under_floor" db:"under_floor"`
	Platform      Platforms    `json:"platform" db:"platform"`
	FActive       ActiveStatus `json:"active" db:"active"`
	CreatedAt     *time.Time   `json:"created_at"  db:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at" db:"updated_at"`
}

// ID returns ID
func (*Publisher) ID() int64 {
	panic("implement me")
}

// FloorCPM returns FloorCPM
func (*Publisher) FloorCPM() int64 {
	panic("implement me")
}

// SoftFloorCPM returns SoftFloorCPM
func (*Publisher) SoftFloorCPM() int64 {
	panic("implement me")
}

// Name returns Name
func (*Publisher) Name() string {
	panic("implement me")
}

// Active returns Active
func (*Publisher) Active() bool {
	panic("implement me")
}

// AcceptedTarget returns AcceptedTarget
func (*Publisher) AcceptedTarget() entity.Target {
	panic("implement me")
}

// Attributes returns Attributes
func (*Publisher) Attributes() interface{} {
	panic("implement me")
}

// BIDType returns BIDType
func (*Publisher) BIDType() entity.BIDType {
	panic("implement me")
}

// MinCPC returns MinCPC
func (*Publisher) MinCPC() int64 {
	panic("implement me")
}

// AcceptedTypes returns AcceptedTypes
func (*Publisher) AcceptedTypes() []entity.AdType {
	panic("implement me")
}

// UnderFloor returns UnderFloor
func (*Publisher) UnderFloor() bool {
	panic("implement me")
}

// Supplier returns Supplier
func (*Publisher) Supplier() entity.Supplier {
	panic("implement me")
}
