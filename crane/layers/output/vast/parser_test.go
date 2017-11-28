package vast

import (
	"net"
	"net/http"
	"testing"

	"time"

	"encoding/xml"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/random"
	"github.com/rs/vmap"
	"github.com/smartystreets/goconvey/convey"
)

type clickProvider struct {
	url string
}

func (c *clickProvider) ClickURL(entity.Slot, entity.Context) string {
	return c.url
}

type imp struct {
	slots    []entity.Slot
	trackID  string
	protocol string
}

type slot struct {
	attr    map[string]interface{}
	winner  entity.Advertise
	showUrl string
	id      string
}

func (s *slot) ID() string {
	return s.id
}

type ad struct {
	width  int
	height int
	id     string
}

func (a *ad) ID() string {
	return a.id
}

func (*ad) Type() entity.AdType {
	panic("implement me")
}

func (*ad) Campaign() entity.Campaign {
	panic("implement me")
}

func (*ad) SetCPM(int64) {
	panic("implement me")
}

func (*ad) CPM() int64 {
	panic("implement me")
}

func (*ad) SetWinnerBID(int64) {
	panic("implement me")
}

func (*ad) WinnerBID() int64 {
	panic("implement me")
}

func (*ad) AdCTR() float64 {
	panic("implement me")
}

func (*ad) SetCTR(float64) {
	panic("implement me")
}

func (*ad) CTR() float64 {
	panic("implement me")
}

func (a *ad) Width() int {
	return a.width
}

func (a *ad) Height() int {
	return a.height
}

func (*ad) Capping() entity.Capping {
	panic("implement me")
}

func (*ad) SetCapping(entity.Capping) {
	panic("implement me")
}

func (*ad) Attributes() map[string]interface{} {
	panic("implement me")
}

func (*ad) Duplicate() entity.Advertise {
	panic("implement me")
}

func (*ad) Media() string {
	panic("implement me")
}

func (*ad) TargetURL() string {
	panic("implement me")
}

func (*slot) TrackID() string {
	panic("implement me")
}

func (*slot) Width() int {
	panic("implement me")
}

func (*slot) Height() int {
	panic("implement me")
}

func (*slot) SetSlotCTR(float64) {
	panic("implement me")
}

func (*slot) SlotCTR() float64 {
	panic("implement me")
}

func (*slot) SetWinnerAdvertise(entity.Advertise) {
	panic("implement me")
}

func (s *slot) WinnerAdvertise() entity.Advertise {
	return s.winner
}

func (*slot) SetShowURL(string) {
	panic("implement me")
}

func (s *slot) ShowURL() string {
	return s.showUrl
}

func (*slot) IsSizeAllowed(int, int) bool {
	panic("implement me")
}

func (s *slot) Attribute() map[string]interface{} {
	return s.attr
}

func (*imp) Request() *http.Request {
	panic("implement me")
}

func (i *imp) TrackID() string {
	return i.trackID
}

func (*imp) ClientID() string {
	panic("implement me")
}

func (*imp) IP() net.IP {
	panic("implement me")
}

func (*imp) UserAgent() string {
	panic("implement me")
}

func (*imp) Publisher() entity.Publisher {
	panic("implement me")
}

func (*imp) Location() entity.Location {
	panic("implement me")
}

func (*imp) OS() entity.OS {
	panic("implement me")
}

func (i *imp) Slots() []entity.Slot {
	return i.slots
}

func (*imp) Category() []entity.Category {
	panic("implement me")
}

func (*imp) Attributes() map[string]string {
	panic("implement me")
}

func (i *imp) Protocol() string {
	return i.protocol
}

func newImpression(iTrackID string, adID string, attr map[string]interface{}, width, height int) entity.Context {
	s := make([]entity.Slot, 0)

	s = append(s, &slot{
		attr: attr,
		winner: &ad{
			width:  width,
			height: height,
			id:     adID,
		},
	})
	a := &imp{
		slots:   s,
		trackID: iTrackID,
	}
	return a
}

func TestParser(t *testing.T) {
	a := data{}
	convey.Convey("sample 1 slot linear", t, func() {
		v1 := vmap.VMAP{}
		attr1 := map[string]interface{}{
			"vast": &slotVastAttribute{
				Duration:  2 * time.Second,
				Offset:    1 * time.Second,
				BreakType: linearType,
			},
		}
		impressionTrackID1 := <-random.ID
		adID1 := <-random.ID

		i := newImpression(impressionTrackID1, adID1, attr1, 300, 250)
		clp := &clickProvider{
			url: "sample click url",
		}
		a.parse(i, clp)
		xml.Unmarshal(a.data, &v1)
		convey.So(1*time.Second, convey.ShouldEqual, *v1.AdBreaks[0].TimeOffset.Duration)
		convey.So(2*time.Second, convey.ShouldEqual, v1.AdBreaks[0].RepeatAfter)
		convey.So(impressionTrackID1, convey.ShouldEqual, v1.AdBreaks[0].BreakID)
		convey.So("linear", convey.ShouldEqual, v1.AdBreaks[0].BreakType)
		convey.So(adID1, convey.ShouldEqual, v1.AdBreaks[0].AdSource.ID)
	})
}
