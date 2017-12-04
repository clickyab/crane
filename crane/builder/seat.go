package builder

import (
	"fmt"
	"strings"

	"errors"

	"strconv"

	"clickyab.com/crane/crane/builder/cynative"
	"clickyab.com/crane/crane/builder/cyslot"
	"clickyab.com/crane/crane/builder/cyvast"
	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/random"
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

// Slot is the slot ID
type Slot struct {
	FPublicID string
	FSize     int
	Type      SlotType

	// App related
	FullScreen string

	// Vast related
	StartTime string
	SkipTime  string

	ExtraParam map[string]string

	ctr         float64
	winnerAd    entity.Advertise
	reserveHash string
	showURL     string
	click       string
}

func (s *Slot) PublicID() string {
	return s.FPublicID
}

func (s *Slot) ReservedHash() string {
	if s.reserveHash == "" {
		s.reserveHash = <-random.ID
	}

	return s.reserveHash
}

func (s *Slot) Width() int {
	w, _ := cyslot.GetSizeByNum(s.FSize)
	width, err := strconv.ParseInt(w, 10, 0)
	assert.Nil(err)
	return int(width)
}

func (s *Slot) Height() int {
	h, _ := cyslot.GetSizeByNum(s.FSize)
	height, err := strconv.ParseInt(h, 10, 0)
	assert.Nil(err)
	return int(height)
}

func (s *Slot) Size() int {
	return s.FSize
}

func (s *Slot) SetSlotCTR(ctr float64) {
	s.ctr = ctr
}

func (s *Slot) SlotCTR() float64 {
	return s.ctr
}

func (s *Slot) SetWinnerAdvertise(wa entity.Advertise) {
	s.winnerAd = wa
}

func (s *Slot) WinnerAdvertise() entity.Advertise {
	return s.winnerAd
}

func (s *Slot) SetShowURL(su string) {
	s.showURL = su
}

func (s *Slot) ShowURL() string {
	return s.showURL
}

func (s *Slot) SetClickURL(c string) {
	s.click = c
}

func (s *Slot) ClickURL() string {
	return s.click
}

func (s *Slot) IsSizeAllowed(w, h int) bool {
	if s.Type == SlotTypeNative {
		return true
	}

	adSize, err := cyslot.GetSize(fmt.Sprintf("%dx%d", w, h))
	if err != nil {
		return false
	}
	return adSize == s.Size()
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
