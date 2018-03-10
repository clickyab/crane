package supplier

import (
	"clickyab.com/crane/demand/entity"
)

type dummy struct {
}

func (s *dummy) SoftFloorCPC(string, string) int64 {
	panic("implement me")
}

func (s *dummy) SoftFloorCPM(string, string) int64 {
	panic("implement me")
}

func (s *dummy) Share() int {
	panic("implement me")
}

func (s *dummy) UnderFloor() bool {
	panic("implement me")
}

func (s *dummy) Name() string {
	return "clickyab"
}

func (s *dummy) Token() string {
	panic("implement me")
}

func (s *dummy) DefaultSoftFloorCPM() int64 {
	panic("implement me")
}

func (s *dummy) DefaultMinBid() int64 {
	panic("implement me")
}
func (s *dummy) Strategy() entity.Strategy {
	panic("implement me")
}
func (s *dummy) DefaultCTR(string, string) float64 {
	panic("implement me")
}

func (s *dummy) AllowCreate() bool {
	return false
}

func (s *dummy) TinyMark() bool {
	panic("implement me")
}

func (s *dummy) TinyLogo() string {
	panic("implement me")
}

func (s *dummy) TinyURL() string {
	panic("implement me")
}

func (s *dummy) ShowDomain() string {
	panic("implement me")
}

func (s *dummy) UserID() int64 {
	panic("implement me")
}

func (s *dummy) Rate() int {
	panic("implement me")
}

// NewClickyab return a dummy clickyab supplier
func NewClickyab() entity.Supplier {
	return &dummy{}
}
