package entities

import (
	"context"
	"fmt"
	"time"

	"database/sql/driver"

	"encoding/gob"
	"io"

	"github.com/clickyab/services/array"
	"github.com/clickyab/services/gettext/t9e"
	"github.com/clickyab/services/kv"
	"github.com/sirupsen/logrus"
)

// StaticSeatType seat type static (banner/native/vast)
type StaticSeatType string

const (

	// NativeStaticSeatType native ad type
	NativeStaticSeatType StaticSeatType = "native"
	// BannerStaticSeatType banner ad type
	BannerStaticSeatType StaticSeatType = "banner"
	// VastStaticSeatType vast ad type
	VastStaticSeatType StaticSeatType = "vast"
)

// IsValid try to validate enum value on ths type
func (e StaticSeatType) IsValid() bool {
	return array.StringInArray(
		string(e),
		string(NativeStaticSeatType),
		string(BannerStaticSeatType),
		string(VastStaticSeatType),
	)
}

// Scan convert the json array ino string slice
func (e *StaticSeatType) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return t9e.G("unsupported type")
	}
	if !StaticSeatType(b).IsValid() {
		return t9e.G("invalid value")
	}
	*e = StaticSeatType(b)
	return nil
}

// Value try to get the string slice representation in database
func (e StaticSeatType) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, t9e.G("invalid status")
	}
	return string(e), nil
}

// StaticSeat model
type StaticSeat struct {
	staticSeat
}

// Publisher return Publisher
func (a *StaticSeat) Publisher() string {
	return a.staticSeat.Publisher
}

// Supplier return Supplier
func (a *StaticSeat) Supplier() string {
	return a.staticSeat.Supplier
}

// Type return Type
func (a *StaticSeat) Type() string {
	return string(a.staticSeat.Type)
}

// Position return Position
func (a *StaticSeat) Position() string {
	return a.staticSeat.Position
}

// From return From
func (a *StaticSeat) From() time.Time {
	return a.staticSeat.From
}

// To return To
func (a *StaticSeat) To() time.Time {
	return a.staticSeat.To
}

// RTBMarkup return RTBMarkup
func (a *StaticSeat) RTBMarkup() string {
	return a.staticSeat.RTBMarkup
}

// Chance return Chance
func (a *StaticSeat) Chance() int {
	return a.staticSeat.Chance
}

// ID return ID
func (a *StaticSeat) ID() int64 {
	return a.staticSeat.ID
}

// Encode is the encode function for serialize object in io writer
func (a *StaticSeat) Encode(w io.Writer) error {
	g := gob.NewEncoder(w)
	return g.Encode(a.staticSeat)
}

// Decode try to decode object from io reader
func (a *StaticSeat) Decode(r io.Reader) error {
	g := gob.NewDecoder(r)
	return g.Decode(a.staticSeat)
}

type staticSeat struct {
	ID        int64          `json:"id" db:"id"`
	Publisher string         `json:"publisher" db:"publisher"`
	Supplier  string         `json:"supplier" db:"supplier"`
	Type      StaticSeatType `json:"type" db:"type"`
	Position  string         `json:"position" db:"position"`
	From      time.Time      `json:"from" db:"from"`
	To        time.Time      `json:"to" db:"to"`
	RTBMarkup string         `json:"rtb_markup" db:"rtb_markup"`
	Chance    int            `json:"chance" db:"chance"`
}

// StaticSeatLoader is the loader of static ads
func StaticSeatLoader(_ context.Context) (map[string]kv.Serializable, error) {
	var res []staticSeat
	t := time.Now()

	query := fmt.Sprintf("SELECT id,publisher,supplier,type,position,`from`,`to`,rtb_markup,chance FROM static_seats WHERE `from` <= ? AND `to` >=?")

	_, err := NewManager().GetRDbMap().Select(
		&res,
		query,
		t,
		t,
	)
	if err != nil {
		return nil, err
	}
	ads := make(map[string]kv.Serializable)
	for i := range res {
		ads[res[i].Publisher+"/"+res[i].Supplier+"/"+string(res[i].Type)+"/"+res[i].Position] = &StaticSeat{staticSeat: res[i]}
	}
	logrus.Debugf("Load %d static ads", len(ads))
	return ads, nil
}
