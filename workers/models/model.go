package models

import (
	"net"
	"time"

	"clickyab.com/crane/demand/entity"
)

// Seat is model for show and click job
type Seat struct {
	AdID         int64              `json:"ad"`
	AdSize       int                `json:"size"`
	SlotPublicID string             `json:"slot"`
	ReserveHash  string             `json:"rh"`
	WinnerBID    float64            `json:"wb"`
	CPM          float64            `json:"cpm"`
	SCPM         float64            `json:"scpm"`
	Type         entity.RequestType `json:"t"`
}

// Impression is model for show and click job
type Impression struct {
	ID                  int64                `json:"id"`
	SeatID              int64                `json:"seat_id"`
	PublisherPageID     int64                `json:"publisher_page_id"`
	CreativesLocationID int64                `json:"creatives_location_id"`
	IP                  net.IP               `json:"ip"`
	CopID               string               `json:"cop"`
	UserAgent           string               `json:"ua"`
	Suspicious          int                  `json:"sp"`
	Referrer            string               `json:"r"`
	ParentURL           string               `json:"par"`
	PublisherID         int64                `json:"pub_id"`
	Publisher           string               `json:"pub"`
	Supplier            string               `json:"sub"`
	Timestamp           time.Time            `json:"ts"`
	PublisherType       entity.PublisherType `json:"pt"`
}
