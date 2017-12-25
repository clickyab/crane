package models

import (
	"net"
	"time"

	"clickyab.com/crane/demand/entity"
)

// Seat is model for show and click job
type Seat struct {
	AdID         int64   `json:"ad"`
	AdSize       int     `json:"size"`
	SlotPublicID string  `json:"slot"`
	ReserveHash  string  `json:"rh"`
	WinnerBID    float64 `json:"wb"`
	CPM          float64 `json:"cpm"`
	SCPM         float64 `json:"scpm"`
}

// Impression is model for show and click job
type Impression struct {
	IP         net.IP             `json:"ip"`
	CopID      string             `json:"cop"`
	UserAgent  string             `json:"ua"`
	Suspicious int                `json:"sp"`
	Referrer   string             `json:"r"`
	ParentURL  string             `json:"par"`
	Publisher  string             `json:"pub"`
	Supplier   string             `json:"sub"`
	Type       entity.RequestType `json:"t"`
	Alexa      bool               `json:"a"`
	Timestamp  time.Time          `json:"ts"`
}
