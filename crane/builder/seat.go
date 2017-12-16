package builder

import (
	"crypto/md5"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"time"

	"net"

	"clickyab.com/crane/crane/builder/internal/cyslot"
	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/framework/router"
	"github.com/clickyab/services/random"
	"github.com/clickyab/services/store/jwt"
)

var (
	showExpire  = config.RegisterDuration("crane.builder.seat.show_exp", 1*time.Hour, "determine how long show url is valid")
	clickExpire = config.RegisterDuration("crane.builder.seat.click_exp", 72*time.Hour, "determine how long click url is valid")
)

// SlotType is the type of slot
type SlotType int

const (
	// SlotTypeWeb web slot
	SlotTypeWeb SlotType = iota
	// SlotTypeApp App slot
	SlotTypeApp
	// SlotTypeVast vast slot
	SlotTypeVast
	// SlotTypeNative native slot
	SlotTypeNative
	// SlotTypeDemand demand slot
	SlotTypeDemand
)

// seat is the seat for input request
type seat struct {
	winnerAd    entity.Advertise
	reserveHash string
	bid         float64
	cpm         float64
	click       string
	show        string

	alexa    bool
	mobile   bool
	iran     bool
	ftype    entity.RequestType
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

	showT  int
	minBid float64

	impTime time.Time

	showExtraParam map[string]string
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

func (s *seat) Publisher() entity.Publisher {
	return s.publisher
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

func (s *seat) SetWinnerAdvertise(wa entity.Advertise, bid float64, cpm float64) {
	s.winnerAd = wa
	s.bid = bid
	s.cpm = cpm
}

func (s *seat) WinnerAdvertise() entity.Advertise {
	return s.winnerAd
}

func (s *seat) ShowURL() string {
	if s.show != "" {
		return s.show
	}
	if s.winnerAd == nil {
		panic("no winner")
	}
	m := md5.New()
	_, _ = m.Write(
		[]byte(s.ReservedHash() + fmt.Sprint(s.size) + s.Type() + s.ua + s.ip.String()),
	)
	data := fmt.Sprintf("%x", m.Sum(nil))
	j := jwt.NewJWT().Encode(map[string]string{
		"aid":  fmt.Sprint(s.winnerAd.ID()),
		"dom":  s.publisher.Name(),
		"sup":  s.publisher.Supplier().Name(),
		"bid":  fmt.Sprint(s.bid),
		"uaip": string(data),
		"pid":  s.publicID,
		"susp": fmt.Sprint(s.susp),
		"now":  fmt.Sprint(time.Now().Unix()),
	}, showExpire.Duration())
	s.winnerAd.ID()
	res := router.MustPath("banner", map[string]string{"rh": s.ReservedHash(), "size": fmt.Sprint(s.size), "jt": j, "type": s.Type()})
	u := url.URL{
		Host:   s.host,
		Scheme: s.protocol.String(),
		Path:   res,
	}

	v := url.Values{}
	v.Set("tid", s.tid)
	v.Set("ref", s.ref)
	v.Set("parent", s.parent)
	for i := range s.showExtraParam {
		v.Set(i, s.showExtraParam[i])
	}
	u.RawQuery = v.Encode()
	s.show = u.String()
	return s.show
}

func (s *seat) ClickURL() string {
	if s.click != "" {
		return s.click
	}
	if s.winnerAd == nil {
		panic("no winner")
	}
	m := md5.New()
	_, _ = m.Write(
		[]byte(s.ReservedHash() + fmt.Sprint(s.size) + s.Type() + s.ua + s.ip.String()),
	)
	data := fmt.Sprintf("%x", m.Sum(nil))
	j := jwt.NewJWT().Encode(map[string]string{
		"aid":  fmt.Sprint(s.winnerAd.ID()),
		"dom":  s.publisher.Name(),
		"sup":  s.publisher.Supplier().Name(),
		"bid":  fmt.Sprint(s.bid),
		"uaip": string(data),
		"susp": fmt.Sprint(s.susp),
		"pid":  s.PublicID(),
		"now":  fmt.Sprint(time.Now().Unix()),
	}, clickExpire.Duration())
	s.winnerAd.ID()
	res, err := router.Path("click", map[string]string{"jt": j, "rh": s.ReservedHash(), "size": fmt.Sprint(s.Size()), "type": s.Type()})
	assert.Nil(err)
	u := url.URL{
		Host:   s.host,
		Scheme: s.protocol.String(),
		Path:   res,
	}
	v := url.Values{}
	v.Set("tid", s.tid)
	v.Set("ref", s.ref)
	v.Set("parent", s.parent)
	u.RawQuery = v.Encode()
	s.click = u.String()
	return s.click
}

func (s *seat) Type() string {
	return string(s.ftype)
}
