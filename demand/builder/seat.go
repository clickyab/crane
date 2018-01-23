package builder

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"time"

	"net"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/internal/cyslot"
	"clickyab.com/crane/demand/internal/hash"
	"github.com/clickyab/services/array"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/store/jwt"
)

var (
	showExpire  = config.RegisterDuration("crane.context.seat.show_exp", 1*time.Hour, "determine how long show url is valid")
	clickExpire = config.RegisterDuration("crane.context.seat.click_exp", 72*time.Hour, "determine how long click url is valid")
)

// seat is the seat for input request
type seat struct {
	winnerAd    entity.Creative
	reserveHash string
	bid         float64
	cpm         float64

	// url cache
	click *url.URL
	imp   *url.URL
	win   *url.URL

	alexa    bool
	mobile   bool
	iran     bool
	ftype    entity.RequestType
	subType  entity.RequestType
	size     int
	publicID string
	ua       string
	ip       net.IP
	tid      string
	ref      string
	parent   string
	susp     int
	protocol entity.Protocol
	// Host return the target host which is different form request.Host and will be used for routing in click, show, etc
	// for example if current request.Host is a.clickyab.com and we want to click url hit b.clickyab.com then Host
	// return b.clickyab.com
	host             string
	minBidPercentage int64
	rate             float64

	publisher entity.Publisher
	ctr       float64

	showT     int
	fatFinger bool
	minBid    float64

	impTime time.Time
	scpm    float64

	// Video related stuff
	mimes []string
}

func (s *seat) FatFinger() bool {
	return s.fatFinger
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
	return int64(math.Ceil(s.minBid*s.rate)/100) * s.minBidPercentage
}

func (s *seat) ShowT() bool {
	if s.showT == 0 {
		if s.mobile && s.iran && s.alexa && rand.Intn(chanceShowT.Int()) == 1 {
			s.showT = 3
		} else {
			s.showT = 2
		}
	}

	return s.showT == 3
}

func (s *seat) CTR() float64 {

	return s.ctr
}

func (s *seat) Width() int {
	w, _ := cyslot.GetSizeByNum(s.size)
	return w

}

func (s *seat) Height() int {
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

func (s *seat) Size() int {
	return s.size
}

func (s *seat) SetWinnerAdvertise(wa entity.Creative, bid float64, cpm float64) {
	s.winnerAd = wa
	s.bid = bid
	s.cpm = cpm
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
		map[string]string{"rh": s.ReservedHash(), "size": fmt.Sprint(s.size), "type": s.Type(), "subtype": s.SubType()},
		s.cpm,
		showExpire.Duration(),
	)
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
		map[string]string{"rh": s.ReservedHash(), "size": fmt.Sprint(s.Size()), "type": s.Type(), "subtype": s.SubType()},
		cpm,
		clickExpire.Duration(),
	)
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
		map[string]string{"rh": s.ReservedHash(), "size": fmt.Sprint(s.Size()), "type": s.Type(), "subtype": s.SubType()},
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
	if s.publisher.Type() == entity.PublisherTypeApp {
		mode = 1
	}
	data := hash.Sign(mode, s.ReservedHash(), fmt.Sprint(s.size), s.Type(), s.ua, s.ip.String())
	ff := "F"
	if s.FatFinger() {
		ff = "T"
	}
	j := jwt.NewJWT().Encode(map[string]string{
		"aid":  fmt.Sprint(s.winnerAd.ID()),
		"dom":  s.publisher.Name(),
		"sup":  s.publisher.Supplier().Name(),
		"bid":  fmt.Sprint(s.bid),
		"uaip": string(data),
		"pid":  s.publicID,
		"susp": fmt.Sprint(s.susp),
		"now":  fmt.Sprint(time.Now().Unix()),
		"cpm":  fmt.Sprint(cpm),
		"ff":   ff,
	}, expire)
	s.winnerAd.ID()
	params["jt"] = j
	res := router.MustPath(
		route,
		params,
	)
	u := &url.URL{
		Host:   s.host,
		Scheme: s.protocol.String(),
		Path:   res,
	}

	v := url.Values{}
	v.Set("tid", s.tid)
	v.Set("ref", s.ref)
	v.Set("parent", s.parent)
	u.RawQuery = v.Encode()
	return u
}

func (s *seat) Type() string {
	return string(s.ftype)
}

func (s *seat) SubType() string {
	return string(s.subType)
}

func (s *seat) Acceptable(advertise entity.Creative) bool {
	if len(s.mimes) > 0 {
		return array.StringInArray(advertise.MimeType(), s.mimes...)
	}
	return true
}
