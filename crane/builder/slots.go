package builder

import (
	"fmt"
	"strings"

	"errors"

	"clickyab.com/gad/builder/cyslot"
	"clickyab.com/gad/builder/cyvast"
	"clickyab.com/gad/utils"
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
	ID       int64
	PublicID string
	Size     int
	Type     SlotType

	// App related
	FullScreen string

	// Vast related
	StartTime string
	SkipTime  string

	ExtraParam map[string]string

	CTR float64

	reserveHash string
	url         string
	click       string
}

// ReserveHash is a simple hash for generating link and waiting to show
func (s *Slot) ReserveHash() string {
	if s.reserveHash == "" {
		s.reserveHash = <-random.ID
	}

	return s.reserveHash
}

// SetURL is the simple url setter for this slot on reserve
func (s *Slot) SetURL(u string) {
	s.url = u
}

// URL is the url getter
func (s *Slot) URL() string {
	return s.url
}

// SizeString return the width and height in string
func (s *Slot) SizeString() (string, string) {
	return utils.GetSizeByNum(s.Size)
}

func (s *Slot) SetClick(u string) {
	s.click = u
}

func (s *Slot) Click() string {
	return s.click
}

// AddWebSlot try to add a web slot to list
func AddWebSlot(pubID string, size int) ShowOptionSetter {
	//validate slot size
	assert.False(utils.InWebSize(size))
	return func(o *Context) (*Context, error) {
		if o.data.Website == nil {
			return nil, errors.New("website not filled")
		}
		finalRes := Slot{
			PublicID: pubID,
			Size:     size,
			Type:     SlotTypeWeb,
		}
		slotID, err := cyslot.GetWebSlotID(pubID, o.data.Website.WID, size)
		if err != nil {
			return nil, errors.New("cant get web slot")
		}
		finalRes.ID = slotID
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
				Size:     utils.NativeAdSize,
				Type:     SlotTypeNative,
			})
			pubIdSize[pub] = utils.NativeAdSize
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
				Size:      size,
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
			Size:       bs,
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
