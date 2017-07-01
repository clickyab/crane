package local

import "clickyab.com/crane/crane/entity"

// Country Country
func (*Location) Country() entity.Country {
	panic("implement me")
}

// Province Province
func (*Location) Province() entity.Province {
	panic("implement me")
}

// LatLon LatLon
func (*Location) LatLon() entity.LatLon {
	panic("implement me")
}

// FloorCPM FloorCPM
func (*Publisher) FloorCPM() int64 {
	panic("implement me")
}

// SoftFloorCPM SoftFloorCPM
func (*Publisher) SoftFloorCPM() int64 {
	panic("implement me")
}

// UnderFloor UnderFloor
func (*Publisher) UnderFloor() bool {
	panic("implement me")
}

// Name Name
func (*Publisher) Name() string {
	panic("implement me")
}

// AcceptedTarget AcceptedTarget
func (*Publisher) AcceptedTarget() entity.Target {
	panic("implement me")
}

// Attributes Attributes
func (*Publisher) Attributes() map[string]interface{} {
	panic("implement me")
}

// BIDType BIDType
func (*Publisher) BIDType() entity.BIDType {
	panic("implement me")
}

// MinCPC MinCPC
func (*Publisher) MinCPC() int64 {
	panic("implement me")
}

// AcceptedTypes AcceptedTypes
func (*Publisher) AcceptedTypes() []entity.AdType {
	panic("implement me")
}

// Supplier Supplier
func (*Publisher) Supplier() string {
	panic("implement me")
}

// ID ID
func (*Slot) ID() string {
	panic("implement me")
}

// TrackID TrackID
func (*Slot) TrackID() string {
	panic("implement me")
}

// Width Width
func (*Slot) Width() int {
	panic("implement me")
}

// Height Height
func (*Slot) Height() int {
	panic("implement me")
}

// SetSlotCTR SetSlotCTR
func (*Slot) SetSlotCTR(float64) {
	panic("implement me")
}

// SlotCTR SlotCTR
func (*Slot) SlotCTR() float64 {
	panic("implement me")
}

// SetWinnerAdvertise SetWinnerAdvertise
func (*Slot) SetWinnerAdvertise(entity.Advertise) {
	panic("implement me")
}

// WinnerAdvertise WinnerAdvertise
func (*Slot) WinnerAdvertise() entity.Advertise {
	panic("implement me")
}

// SetShowURL SetShowURL
func (*Slot) SetShowURL(string) {
	panic("implement me")
}

// ShowURL ShowURL
func (*Slot) ShowURL() string {
	panic("implement me")
}

// IsSizeAllowed IsSizeAllowed
func (*Slot) IsSizeAllowed(int, int) bool {
	panic("implement me")
}

// Attribute Attribute
func (*Slot) Attribute() map[string]interface{} {
	panic("implement me")
}

// SetAdvertise SetAdvertise
func (*Slot) SetAdvertise(a entity.Advertise) {
	panic("implement me")
}

// Advertise Advertise
func (*Slot) Advertise() entity.Advertise {
	panic("implement me")
}
