package rtb

import (
	"context"
	"net"
	"net/http"
	"octopus/exchange"
	"services/random"
	"strconv"
	"testing"

	x "github.com/smartystreets/assertions"


	. "github.com/smartystreets/goconvey/convey"
)

type situation struct {
	Case          int
	SCPMFloor     int64
	CPMFloor      int64
	MarginPercent int
	DSCPMFloor    int64
	DCPMFloor     int64
	UnderFloor    bool
	Demands       []dem
	WinnerDemand  int64
	Profit3rdAd   float64
	ProfitSuplier float64
	description   string
	Expectation   func(interface{}, ...interface{}) string
}

type dem struct {
	MaxCPM   int64
	HandyCap int64
}

func cases() []situation {
	return []situation{
		{1, 300, 200, 10, 330, 220, false, []dem{{470, 100}, {440, 110}}, 440, 44, 396, "", x.ShouldEqual},
		{2, 300, 200, 10, 330, 220, true, []dem{{470, 100}, {440, 110}}, 440, 44, 396, "", x.ShouldEqual},
		{3, 300, 200, 10, 330, 220, false, []dem{{470, 100}, {440, 100}}, 441, 44.1, 396.9, "", x.ShouldEqual},
		{4, 300, 200, 10, 330, 220, true, []dem{{470, 100}, {440, 100}}, 441, 44.1, 396.9, "", x.ShouldEqual},
		{5, 300, 200, 10, 330, 220, false, []dem{{470, 120}, {230, 115}}, 331, 33, 297, "", x.ShouldEqual},
		{6, 300, 200, 10, 330, 220, true, []dem{{470, 120}, {230, 115}}, 331, 33, 297, "", x.ShouldEqual},
		{7, 300, 200, 10, 330, 220, false, []dem{{340, 100}, {300, 125}}, 300, 30, 270, "", x.ShouldEqual},
		{8, 300, 200, 10, 330, 220, true, []dem{{340, 100}, {300, 125}}, 300, 30, 270, "", x.ShouldEqual},
		{9, 300, 200, 10, 330, 220, false, []dem{{230, 100}, {250, 110}}, 231, 23.1, 207.9, "", x.ShouldEqual},
		{10, 300, 200, 10, 330, 220, true, []dem{{230, 100}, {250, 110}}, 231, 23.1, 207.9, "", x.ShouldEqual},
		{11, 300, 200, 10, 330, 220, false, []dem{{200, 120}, {230, 90}}, 200, 22, 198, "", x.ShouldEqual},
		{12, 300, 200, 10, 330, 220, true, []dem{{200, 120}, {230, 90}}, 200, 20, 180, "", x.ShouldEqual},
		{14, 300, 200, 10, 330, 220, true, []dem{{190, 100}, {185, 105}}, 185, 18.5, 166.5, "", x.ShouldEqual},
		{15, 300, 200, 10, 330, 220, false, []dem{{250, 100}, {250, 100}}, 250, 25, 225, "", x.ShouldEqual},
		{16, 300, 200, 10, 330, 220, false, []dem{{350, 100}}, 331, 33, 297, "", x.ShouldEqual},
		{17, 300, 200, 10, 330, 220, true, []dem{{350, 100}}, 331, 33, 297, "", x.ShouldEqual},
		{18, 300, 200, 10, 330, 220, false, []dem{{310, 100}}, 221, 22, 198, "", x.ShouldEqual},
		{19, 300, 200, 10, 330, 220, true, []dem{{310, 100}}, 221, 22, 198, "", x.ShouldEqual},
		{20, 300, 200, 10, 330, 220, false, []dem{{240, 120}}, 221, 22, 198, "", x.ShouldEqual},
		{21, 300, 200, 10, 330, 220, true, []dem{{240, 120}}, 221, 22, 198, "", x.ShouldEqual},
		{23, 300, 200, 10, 330, 220, true, []dem{{210, 100}}, 210, 21, 189, "", x.ShouldEqual},
		{25, 300, 200, 10, 330, 220, true, []dem{{210, 80}}, 210, 21, 189, "", x.ShouldEqual},
		{27, 300, 200, 10, 330, 220, true, []dem{{210, 80}}, 210, 21, 189, "", x.ShouldEqual},
		{29, 300, 200, 10, 330, 220, true, []dem{{190, 100}}, 190, 21, 189, "", x.ShouldEqual},
	}

}

type Advertise struct {
	cpm,
	win int64
	demand exchange.Demand
}

func (a *Advertise) ID() string              { panic("Advertise") }
func (a *Advertise) MaxCPM() int64           { return a.cpm }
func (a *Advertise) Width() int              { panic("Advertise") }
func (a *Advertise) Height() int             { panic("Advertise") }
func (a *Advertise) URL() string             { panic("Advertise") }
func (a *Advertise) TrackID() string         { panic("Advertise") }
func (a *Advertise) SetWinnerCPM(w int64)    { a.win = w }
func (a *Advertise) WinnerCPM() int64        { return a.win }
func (a *Advertise) Demand() exchange.Demand { return a.demand }

type Impression struct {
	publisher  exchange.Publisher
	underFloor bool
}

func (i *Impression) TrackID() string                    { panic("Impression") }
func (i *Impression) IP() net.IP                         { panic("Impression") }
func (i *Impression) UserAgent() string                  { panic("Impression") }
func (i *Impression) Source() exchange.Publisher         { return i.publisher }
func (i *Impression) Location() exchange.Location        { panic("Impression") }
func (i *Impression) Attributes() map[string]interface{} { panic("Impression") }
func (i *Impression) Slots() []exchange.Slot             { panic("Impression") }
func (i *Impression) Category() []exchange.Category      { panic("Impression") }
func (i *Impression) Type() exchange.ImpressionType      { panic("Impression") }
func (i *Impression) UnderFloor() bool                   { return i.underFloor }

type Demand struct {
	handicap int64
}

func (d *Demand) Name() string { panic("Demand") }
func (d *Demand) Provide(context.Context, exchange.Impression, chan map[string]exchange.Advertise) {
	panic("Demand")
}
func (d *Demand) Win(context.Context, string, int64)                         { panic("Demand") }
func (d *Demand) Status(context.Context, http.ResponseWriter, *http.Request) { panic("Demand") }
func (d *Demand) Handicap() int64                                            { return d.handicap }
func (d *Demand) CallRate() int                                              { panic("Demand") }

type Publisher struct {
	floorCPM,
	softFloorCPM int64
}

func (p *Publisher) Name() string                       { panic("publisher") }
func (p *Publisher) FloorCPM() int64                    { return p.floorCPM }
func (p *Publisher) SoftFloorCPM() int64                { return p.softFloorCPM }
func (p *Publisher) Attributes() map[string]interface{} { panic("publisher") }
func (p *Publisher) Supplier() exchange.Supplier        { panic("publisher") }

func TestSelect(t *testing.T) {

	for _, u := range cases() {
		Convey("SelectCPM function test with case number "+strconv.Itoa(u.Case), t, func() {
			p := &Publisher{u.DCPMFloor, u.DSCPMFloor}
			m := &Impression{p, u.UnderFloor}

			ads := make([]exchange.Advertise, 0)

			for _, a := range u.Demands {
				d := &Demand{a.HandyCap}
				ad := &Advertise{a.MaxCPM, u.WinnerDemand, d}
				ads = append(ads, ad)
			}

			slots := make(map[string][]exchange.Advertise)
			id := <-random.ID
			slots[id] = ads
			res := SelectCPM(m, slots)
			So(res[id].WinnerCPM(), u.Expectation, u.WinnerDemand)
		})
	}

	Convey("SelectCPM function Should return nil", t, func() {
		m := &Impression{}
		id := <-random.ID
		s := make(map[string][]exchange.Advertise)
		s[id] = nil
		r := SelectCPM(m, s)
		So(r[id], ShouldBeNil)
	})

}
