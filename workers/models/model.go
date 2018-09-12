package models

import (
	"net"
	"time"

	"clickyab.com/crane/demand/entity"
)

// Seat is model for show and click job
type Seat struct {
	AdID         int32              `json:"ad"`
	AdSize       int32              `json:"size"`
	SlotPublicID string             `json:"slot"`
	ReserveHash  string             `json:"rh"`
	WinnerBID    float64            `json:"wb"`
	CPM          float64            `json:"cpm"`
	SCPM         float64            `json:"scpm"`
	Type         entity.RequestType `json:"t"`
}

// Impression is model for show and click job
type Impression struct {
	IP            net.IP               `json:"ip"`
	CopID         string               `json:"cop"`
	UserAgent     string               `json:"ua"`
	Suspicious    int                  `json:"sp"`
	Referrer      string               `json:"r"`
	ParentURL     string               `json:"par"`
	Publisher     string               `json:"pub"`
	Supplier      string               `json:"sub"`
	Timestamp     time.Time            `json:"ts"`
	PublisherType entity.PublisherType `json:"pt"`
}
