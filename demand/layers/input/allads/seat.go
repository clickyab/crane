package allads

import (
	"math"
	"net/url"
	"time"

	"clickyab.com/crane/demand/entity"
	"github.com/clickyab/services/array"
)

type seat struct {
	mimes   []string
	rq      entity.RequestType
	ctr     float64
	minBid  float64
	context entity.Context
	rate    float64
	size    int
}

func (s *seat) PublicID() string {
	panic("implement me")
}

func (s *seat) ReservedHash() string {
	panic("implement me")
}

func (s *seat) Width() int {
	panic("implement me")
}

func (s *seat) Height() int {
	panic("implement me")
}

func (s *seat) Size() int {
	panic("implement me")
}

func (s *seat) Bid() float64 {
	panic("implement me")
}

func (s *seat) CPM() float64 {
	panic("implement me")
}

func (s *seat) SetWinnerAdvertise(entity.Creative, float64, float64) {
	panic("implement me")
}

func (s *seat) WinnerAdvertise() entity.Creative {
	panic("implement me")
}

func (s *seat) ImpressionURL() *url.URL {
	panic("implement me")
}

func (s *seat) ClickURL() *url.URL {
	panic("implement me")
}

func (s *seat) WinNoticeRequest() *url.URL {
	panic("implement me")
}

func (s *seat) CTR() float64 {
	return s.ctr
}

func (s *seat) Type() entity.InputType {
	panic("implement me")
}

func (s *seat) RequestType() entity.RequestType {
	return s.rq
}

func (s *seat) MinBid() int64 {

	return int64(math.Ceil((s.minBid * s.rate / 100) * float64(s.context.MinBIDPercentage())))
}

func (s *seat) ImpressionTime() time.Time {
	panic("implement me")
}

func (s *seat) SupplierCPM() float64 {
	panic("implement me")
}

func (s *seat) FatFinger() bool {
	panic("implement me")
}

func (s *seat) Acceptable(advertise entity.Creative) bool {
	// this function, handle banner only
	if !s.genericTests(advertise) {
		return false
	}

	if s.size != advertise.Size() {
		return false
	}

	switch s.context.Publisher().Type() {
	case entity.PublisherTypeApp:
		if advertise.Type() == entity.AdTypeVideo || advertise.Type() == entity.AdTypeDynamic {
			return false
		}
		if advertise.Target() != entity.TargetApp {
			return false
		}
	case entity.PublisherTypeWeb:
		if advertise.Target() != entity.TargetWeb {
			return advertise.Campaign().Web() || advertise.Campaign().WebMobile()
		}

	default:
		panic("invalid type")
	}

	return true
}

func (s seat) genericTests(advertise entity.Creative) bool {
	// if the seat has mime setting, make sure we honor it.
	if len(s.mimes) > 0 {
		if !array.StringInArray(advertise.MimeType(), s.mimes...) {
			return false
		}
	}

	return true
}
