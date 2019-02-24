package builder

import (
	"fmt"
	"math"
	"net/url"
	"time"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/internal/cyslot"
	"clickyab.com/crane/internal/hash"
	"github.com/clickyab/services/array"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/store/jwt"
)

var (
	showExpire    = config.RegisterDuration("crane.context.seat.show_exp", 1*time.Hour, "determine how long show url is valid")
	clickExpire   = config.RegisterDuration("crane.context.seat.click_exp", 72*time.Hour, "determine how long click url is valid")
	currentRegion = config.RegisterString("crane.regions.current", "fr", "determine current region")
)

// seat is the seat for input request
type seat struct {
	context     *Context
	winnerAd    entity.Creative
	reserveHash string
	bid         float64
	cpm         float64
	did         string

	// url cache
	click *url.URL
	imp   *url.URL
	win   *url.URL

	size     int32
	publicID string

	rate float64

	ctr float32

	minBid  float64
	minCPC  float64
	softCPM float64
	minCPM  float64

	impTime time.Time
	scpm    float64

	// Video related stuff
	mimes       []string
	requestType entity.RequestType
	cpID        int32
	cpAdID      int32
	adID        int32
	cpName      string
	tr          string
}

func (s *seat) TargetURL() string {
	return s.tr
}

func (s *seat) CampaignName() string {
	return s.cpName
}

func (s *seat) CampaignID() int32 {
	return s.cpID
}

func (s *seat) CampaignAdID() int32 {
	return s.cpAdID
}

func (s *seat) AdID() int32 {
	return s.adID
}

func (s *seat) DID() string {
	return s.did
}

func (s *seat) RequestType() entity.RequestType {
	return s.requestType
}

func (s *seat) FatFinger() bool {
	return s.context.FatFinger()
}

func (s *seat) SupplierCPM() float64 {
	return s.scpm
}

func (s *seat) ImpressionTime() time.Time {
	return s.impTime
}

func (s *seat) CPM() float64 {
	return s.cpm
}

// MinBid return the current minimum bid, apply the min bid percentage and
// rate.
func (s *seat) MinBid() int64 {
	return int64(math.Ceil((float64(s.minBid) * s.rate / 100) * float64(s.context.MinBIDPercentage())))
}

func (s *seat) CTR() float32 {

	return s.ctr
}

func (s *seat) Width() int32 {
	w, _ := cyslot.GetSizeByNum(s.size)
	return w

}

func (s *seat) Height() int32 {
	_, h := cyslot.GetSizeByNum(s.size)
	return h
}

func (s *seat) Bid() float64 {
	return s.bid
}

func (s *seat) PublicID() string {
	return s.publicID
}

func (s *seat) ReservedHash() string {
	if s.reserveHash == "" {
		s.reserveHash = <-random.ID
	}
	return s.reserveHash
}

func (s *seat) Size() int32 {
	return s.size
}

func (s *seat) SetWinnerAdvertise(wa entity.Creative, bid float64, cpm float64) {
	s.winnerAd = wa
	s.bid = bid
	s.cpm = decShare(s.context.publisher.Supplier(), cpm)
	s.adID = wa.ID()
	s.cpAdID = wa.CampaignAdID()
	s.cpID = wa.Campaign().ID()
	s.cpName = wa.Campaign().Name()
	s.tr = wa.TargetURL()
}

func (s *seat) WinnerAdvertise() entity.Creative {
	return s.winnerAd
}

func (s *seat) ImpressionURL() *url.URL {
	if s.imp != nil {
		return s.imp
	}
	if s.winnerAd == nil {
		panic("no winner")
	}

	s.imp = s.makeURL(
		"banner",
		map[string]string{
			"rh":      s.ReservedHash(),
			"size":    fmt.Sprint(s.Size()),
			"type":    s.Type().String(),
			"subtype": s.RequestType().String(),
			"pt":      s.context.publisher.Type().String(),
		},
		s.cpm,
		showExpire.Duration(),
	)

	//check for markup if markup is true act like native impression
	// TODO : move to better place :(
	if s.context.bannerMarkup {
		s.imp = s.makeURL(
			"pixel",
			map[string]string{
				"rh":      s.ReservedHash(),
				"size":    fmt.Sprint(s.Size()),
				"type":    s.Type().String(),
				"subtype": s.RequestType().String(),
				"pt":      s.context.publisher.Type().String(),
			},
			s.cpm,
			showExpire.Duration(),
		)
	}
	return s.imp
}

func (s *seat) ClickURL() *url.URL {
	if s.click != nil {
		return s.click
	}
	if s.winnerAd == nil {
		panic("no winner")
	}
	cpm := s.cpm
	if s.scpm != 0 {
		cpm = s.scpm
	}

	s.click = s.makeURL(
		"click",
		map[string]string{
			"rh":      s.ReservedHash(),
			"size":    fmt.Sprint(s.Size()),
			"type":    s.Type().String(),
			"subtype": s.RequestType().String(),
			"pt":      s.context.publisher.Type().String(),
		},
		cpm,
		clickExpire.Duration(),
	)
	// a hack, to prevent https links on click. it seems google analytic has some kind of problem with
	// https referrer : https://www.e-nor.com/blog/google-analytics/https-to-http-secure-to-nonsecure-referrer-loss
	s.click.Scheme = string(entity.HTTPS)

	return s.click
}

func (s *seat) WinNoticeRequest() *url.URL {
	if s.win != nil {
		return s.win
	}
	if s.winnerAd == nil {
		panic("no winner")
	}

	s.win = s.makeURL(
		"notice",
		map[string]string{
			"rh":      s.ReservedHash(),
			"size":    fmt.Sprint(s.Size()),
			"type":    s.Type().String(),
			"subtype": s.RequestType().String(),
			"pt":      s.context.publisher.Type().String(),
		},
		s.cpm,
		time.Hour, // TODO : fix me when there is actually a code to handle it
	)
	return s.win
}

func (s *seat) makeURL(route string, params map[string]string, cpm float64, expire time.Duration) *url.URL {
	if s.winnerAd == nil {
		panic("no winner")
	}
	mode := 0
	if s.context.Publisher().Type() == entity.PublisherTypeApp {
		mode = 1
	}
	data := hash.Sign(mode, s.ReservedHash(), fmt.Sprint(s.size), s.context.Type().String(), s.context.UserAgent(), s.context.IP().String())
	ff := "F"
	if s.FatFinger() {
		ff = "T"
	}
	tiny := "F"
	if s.context.Tiny() {
		tiny = "T"
	}
	j := jwt.NewJWT().Encode(map[string]string{
		"aid":    fmt.Sprint(s.winnerAd.ID()),
		"dom":    s.context.Publisher().Name(),
		"sup":    s.context.Publisher().Supplier().Name(),
		"bid":    fmt.Sprint(s.bid),
		"uaip":   string(data),
		"pid":    s.publicID,
		"susp":   fmt.Sprint(s.context.Suspicious()),
		"now":    fmt.Sprint(time.Now().Unix()),
		"cpm":    fmt.Sprint(cpm),
		"ff":     ff,
		"pt":     s.context.Publisher().Type().String(),
		"t":      tiny,
		"cmode":  fmt.Sprint(s.context.Capping()),
		"did":    s.did,
		"cpid":   fmt.Sprint(s.winnerAd.Campaign().ID()),
		"cpadid": fmt.Sprint(s.winnerAd.CampaignAdID()),
		"tr":     s.winnerAd.TargetURL(),
		"cpn":    s.winnerAd.Campaign().Name(),
	}, expire)
	s.winnerAd.ID()
	params["jt"] = j
	res := router.MustPath(
		route,
		params,
	)
	u := &url.URL{
		Host:   s.context.host,
		Scheme: s.context.Protocol().String(),
		Path:   res,
	}
	v := url.Values{}
	v.Set("uid", s.context.user.ID())
	v.Set("tid", s.context.tid)
	v.Set("ref", s.context.referrer)
	v.Set("parent", s.context.parent)
	v.Set("reg", currentRegion.String())
	u.RawQuery = v.Encode()
	u.Scheme = "https"
	return u
}

func (s *seat) Type() entity.InputType {
	return s.context.Type()
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

// Acceptable play a crucial role here.
func (s *seat) Acceptable(advertise entity.Creative) bool {
	// this function, handle banner only
	if !s.genericTests(advertise) {
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

// MinCPC return min cpc
func (s *seat) MinCPC() float64 {
	return s.minCPC
}

// MinCPM return min cpm
func (s *seat) MinCPM() float64 {
	return s.minCPM
}

// SoftCPM is the soft lower cpm
func (s *seat) SoftCPM() float64 {
	return s.softCPM
}
