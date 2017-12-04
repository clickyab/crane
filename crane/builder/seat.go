package builder

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"strings"
	"time"

	"errors"

	"clickyab.com/crane/crane/builder/cynative"
	"clickyab.com/crane/crane/builder/cyslot"
	"clickyab.com/crane/crane/builder/cyvast"
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
	FPublicID   string
	FSize       int
	ctr         float64
	winnerAd    entity.Advertise
	reserveHash string
	bid         float64
	click       string
	show        string

	publisherID string
	ua          string
	ip          string
	tid         string
	ref         string
	susp        string
	parent      string
	protocol    entity.Protocol
	// Host return the target host which is different form request.Host and will be used for routing in click, show, etc
	// for example if current request.Host is a.clickyab.com and we want to click url hit b.clickyab.com then Host
	// return b.clickyab.com
	host string
}

func (s *seat) Width() int {
	w, _ := cyslot.GetSizeByNum(s.FSize)
	return w

}

func (s *seat) Height() int {
	_, h := cyslot.GetSizeByNum(s.FSize)
	return h
}

func (s *seat) Bid() float64 {
	return s.bid
}

func (s *seat) PublicID() string {
	return s.FPublicID
}

func (s *seat) ReservedHash() string {
	if s.reserveHash == "" {
		s.reserveHash = <-random.ID
	}
	return s.reserveHash
}

func (s *seat) Size() int {
	return s.FSize
}

func (s *seat) SetWinnerAdvertise(wa entity.Advertise, p float64) {
	s.winnerAd = wa
	s.bid = p
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
	j := jwt.NewJWT().Encode(map[string]string{
		"a": fmt.Sprint(s.winnerAd.ID()),
		"b": s.publisherID,
		"c": fmt.Sprint(s.bid),
		"d": s.publisherID,
		"e": string(m.Sum([]byte(s.ua + s.ip))),
	}, showExpire.Duration())
	s.winnerAd.ID()
	res, err := router.Path("show", map[string]string{"j": j, "t": s.tid, "ref": s.ref, "parent": s.parent, "s": fmt.Sprint(s.Size())})
	assert.Nil(err)
	u := url.URL{
		Host:   s.host,
		Scheme: s.protocol.String(),
		Path:   res,
	}
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
	j := jwt.NewJWT().Encode(map[string]string{
		"a": fmt.Sprint(s.winnerAd.ID()),
		"b": s.publisherID,
		"c": fmt.Sprint(s.bid),
		"d": s.publisherID,
		"e": string(m.Sum([]byte(s.ua + s.ip))),
		"f": s.susp,
	}, showExpire.Duration())
	s.winnerAd.ID()
	res, err := router.Path("click", map[string]string{"j": j, "t": s.tid, "ref": s.ref, "parent": s.parent, "s": fmt.Sprint(s.Size())})
	assert.Nil(err)
	u := url.URL{
		Host:   s.host,
		Scheme: s.protocol.String(),
		Path:   res,
	}
	s.click = u.String()
	return s.click
}

// AddWebSlot try to add a web slot to list
func AddWebSlot(pubID string, size int) ShowOptionSetter {
	//validate slot size
	assert.True(cyslot.ValidWebSlotSize(size))
	return func(o *Context) (*Context, error) {
		if o.data.Website == nil {
			return nil, errors.New("website not filled")
		}
		finalRes := Slot{
			PublicID: pubID,
			FSize:    size,
			Type:     SlotTypeWeb,
		}
		slotID, err := cyslot.GetWebSlotID(pubID, o.data.Website.WID, size)
		if err != nil {
			return nil, errors.New("cant get web slot")
		}
		finalRes.FID = slotID
		o.rtb.Slots = append(o.rtb.Slots, &finalRes)
		return o, nil
	}
}

// AddNativeSlot try to add native slot to list
func AddNativeSlot(count int, pubIDBase string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if o.data.Website == nil {
			return nil, errors.New("website not filled")
		}
		var pubIdSize = make(map[string]int)
		for i := 0; i < count; i++ {
			pub := fmt.Sprintf(pubIDBase+"%d", i)
			o.rtb.Slots = append(o.rtb.Slots, &Slot{
				PublicID: pub,
				FSize:    cynative.NativeSlotSize,
				Type:     SlotTypeNative,
			})
			pubIdSize[pub] = cynative.NativeSlotSize
		}
		fSlotRes, err := cyslot.GetCommonSlotIDs(pubIdSize, o.data.Website.WID)
		if err != nil {
			return nil, errors.New("cant get native slot")
		}
		for j := range o.rtb.Slots {
			o.rtb.Slots[j].ID = fSlotRes[o.rtb.Slots[j].PublicID]
		}
		return o, nil
	}
}

// AddVastSlot try to add a new vast slot
func AddVastSlot(basePubID string, l string, first, mid, last bool) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		var pubIdSize = make(map[string]int)
		_, length := cyvast.MakeVastLen(l, first, mid, last)
		var i int
		for m := range length {
			i++
			lenType := length[m][0]
			if lenType != "linear" && o.common.Mobile {
				continue
			}
			pub := fmt.Sprintf("%s%s", basePubID, length[m][1])

			size := cyvast.VastNonLinearSize
			stop := ""
			if lenType == "linear" {
				size = cyvast.VastLinearSize
				stop = length[m][3]
			}
			pubIdSize[pub] = size

			s := &Slot{
				FSize:     size,
				PublicID:  pub,
				Type:      SlotTypeVast,
				StartTime: length[m][2],
				SkipTime:  stop,
				ExtraParam: map[string]string{
					"pos":  m,
					"type": length[m][2],
					"l":    lenType,
				},
			}

			o.rtb.Slots = append(o.rtb.Slots, s)
		}
		fSlotRes, err := cyslot.GetCommonSlotIDs(pubIdSize, o.data.Website.WID)
		if err != nil {
			return nil, errors.New("cant get vast slot")
		}
		for j := range o.rtb.Slots {
			o.rtb.Slots[j].ID = fSlotRes[o.rtb.Slots[j].PublicID]
		}
		return o, nil
	}
}

// AddAppSlot try to add a slot for application
func AddAppSlot(adsMedia string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if o.data.App == nil {
			return nil, fmt.Errorf("no publisher is set")
		}
		var (
			bs   int
			full string
		)
		switch strings.ToLower(adsMedia) {
		case "banner":
			bs = 8
		case "largebanner":
			bs = 3
		case "xlargebannerportrait":
			bs = 16
		case "fullbannerportrait":
			bs = 16
			full = "portrait"
		case "xlargebannerlandscap":
			bs = 17
		case "fullbannerlandscape":
			bs = 17
			full = "landscape"
		default:
			return nil, fmt.Errorf("invalid Size string %s", adsMedia)
		}

		slotString := fmt.Sprintf("%d0%d0%d", o.data.App.ID, o.data.App.UserID, bs)
		fRes := Slot{
			Type:       SlotTypeApp,
			FSize:      bs,
			PublicID:   slotString,
			FullScreen: full,
		}
		slotID, err := cyslot.GetAppSlotID(slotString, o.data.App.ID, bs)
		if err != nil {
			return nil, errors.New("cant get app slot")
		}
		fRes.ID = slotID
		o.rtb.Slots = append(o.rtb.Slots, &fRes)
		return o, nil
	}
}
