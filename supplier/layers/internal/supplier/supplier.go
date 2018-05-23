package supplier

import (
	"clickyab.com/crane/demand/entity"
)

type dummy struct {
}

func (s *dummy) SoftFloorCPC(string, string) int64 {
	panic("implement dummy SoftFloorCPC not valid")
}

func (s *dummy) SoftFloorCPM(string, string) int64 {
	panic("implement dummy SoftFloorCPM not valid")
}

func (s *dummy) Share() int {
	panic("implement dummy Share not valid")
}

func (s *dummy) UnderFloor() bool {
	panic("implement dummy UnderFloor not valid")
}

func (s *dummy) Name() string {
	return "clickyab"
}

func (s *dummy) Token() string {
	panic("implement dummy Token not valid")
}

func (s *dummy) DefaultSoftFloorCPM() int64 {
	panic("implement dummy DefaultSoftFloorCPM not valid")
}

func (s *dummy) DefaultMinBid() int64 {
	panic("implement dummy DefaultMinBid not valid")
}
func (s *dummy) Strategy() entity.Strategy {
	panic("implement dummy Strategy not valid")
}
func (s *dummy) DefaultCTR(string, string) float64 {
	panic("implement dummy DefaultCTR not valid")
}

func (s *dummy) AllowCreate() bool {
	return false
}

func (s *dummy) TinyMark() bool {
	panic("implement dummy TinyMark not valid")
}

func (s *dummy) TinyLogo() string {
	panic("implement dummy TinyLogo not valid")
}

func (s *dummy) TinyURL() string {
	panic("implement dummy TinyURL not valid")
}

func (s *dummy) ShowDomain() string {
	panic("implement dummy ShowDomain not valid")
}

func (s *dummy) UserID() int64 {
	panic("implement dummy UserID not valid")
}

func (s *dummy) Rate() int {
	panic("implement dummy Rate not valid")
}

// NewClickyab return a dummy clickyab supplier
func NewClickyab() entity.Supplier {
	return &dummy{}
}
