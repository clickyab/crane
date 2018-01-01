package app

import (
	"clickyab.com/crane/demand/entity"
)

type supplier struct {
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

func (s *supplier) DefaultFloorCPM() int64 {
	panic("implement me")
}

func (s *supplier) DefaultSoftFloorCPM() int64 {
	panic("implement me")
}

func (s *supplier) DefaultMinBid() int64 {
	panic("implement me")
}

func (s *supplier) BidType() entity.BIDType {
	panic("implement me")
}

func (s *supplier) DefaultCTR() float64 {
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
