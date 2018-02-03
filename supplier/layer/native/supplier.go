package native

import "clickyab.com/crane/demand/entity"

type supplier struct {
}

func (*supplier) SoftFloorCPM(string, string) int64 {
	panic("implement me")
}

func (*supplier) Name() string {
	return "clickyab"
}

func (*supplier) Token() string {
	panic("implement me")
}

func (*supplier) DefaultFloorCPM() int64 {
	panic("implement me")
}

func (*supplier) DefaultSoftFloorCPM() int64 {
	panic("implement me")
}

func (*supplier) DefaultMinBid() int64 {
	panic("implement me")
}

func (*supplier) Strategy() entity.Strategy {
	panic("implement me")
}
func (*supplier) DefaultCTR(string, string) float64 {
	panic("implement me")
}

func (*supplier) AllowCreate() bool {
	panic("implement me")
}

func (*supplier) TinyMark() bool {
	panic("implement me")
}

func (*supplier) TinyLogo() string {
	panic("implement me")
}

func (*supplier) TinyURL() string {
	panic("implement me")
}

func (*supplier) ShowDomain() string {
	panic("implement me")
}

func (*supplier) UserID() int64 {
	panic("implement me")
}

func (*supplier) Rate() int {
	panic("implement me")
}

func (*supplier) UnderFloor() bool {
	panic("implement me")
}

func (*supplier) Share() int {
	panic("implement me")
}
