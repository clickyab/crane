package web

import (
	"clickyab.com/crane/demand/entity"
)

type supplier struct {
}

func (s *supplier) SoftFloorCPM(string, string) int64 {
	panic("implement me")
}

func (s *supplier) UnderFloor() bool {
	panic("implement me")
}

func (s *supplier) Name() string {
	return "clickyab"
}

func (s *supplier) Token() string {
	panic("implement me")
}

func (s *supplier) DefaultSoftFloorCPM() int64 {
	panic("implement me")
}

func (s *supplier) DefaultMinBid() int64 {
	panic("implement me")
}
func (s *supplier) Strategy() entity.Strategy {
	panic("implement me")
}
func (s *supplier) DefaultCTR(string, string) float64 {
	panic("implement me")
}

func (s *supplier) AllowCreate() bool {
	return false
}

func (s *supplier) TinyMark() bool {
	panic("implement me")
}

func (s *supplier) TinyLogo() string {
	panic("implement me")
}

func (s *supplier) TinyURL() string {
	panic("implement me")
}

func (s *supplier) ShowDomain() string {
	panic("implement me")
}

func (s *supplier) UserID() int64 {
	panic("implement me")
}

func (s *supplier) Rate() int {
	panic("implement me")
}

func (s *supplier) Share() int {
	panic("implement me")
}
