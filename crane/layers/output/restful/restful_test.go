package restful

import (
	"bytes"
	"testing"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/entity/mock_entity"

	"regexp"

	"github.com/clickyab/services/random"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
)

func TestRestful(t *testing.T) {
	convey.Convey("Restfull Should have 10 ad", t, func() {

		b := mockImpression(10, t, render{})
		p := regexp.MustCompile("max_cpm")
		res := p.FindAll(b, -1)
		convey.So(len(res), convey.ShouldEqual, 10)
	})

}

func mockImpression(c int, t *testing.T, r entity.Renderer) []byte {
	ct := gomock.NewController(t)
	m := mock_entity.NewMockImpression(ct)
	as := make([]entity.Slot, 0)
	for i := 0; i < c; i++ {
		s := mock_entity.NewMockSlot(ct)

		a := mock_entity.NewMockAdvertise(ct)
		a.EXPECT().WinnerBID().Return(int64(1000)).AnyTimes()
		a.EXPECT().TargetURL().Return("Target URL").AnyTimes()
		a.EXPECT().ID().Return("id").AnyTimes()

		s.EXPECT().SlotCTR().Return(1.).AnyTimes()
		s.EXPECT().WinnerAdvertise().Return(a).AnyTimes()
		s.EXPECT().Width().Return(100).AnyTimes()
		s.EXPECT().Height().Return(100).AnyTimes()
		s.EXPECT().ShowURL().Return("Show URL").AnyTimes()
		s.EXPECT().TrackID().Return(<-random.ID).AnyTimes()

		as = append(as, s)

	}
	m.EXPECT().Slots().Return(as)

	buf := &bytes.Buffer{}

	cp := mock_entity.NewMockClickProvider(ct)
	cp.EXPECT().ClickURL(gomock.Any(), gomock.Any()).Return("FAKE URL")

	r.Render(buf, m, cp)
	return buf.Bytes()
}
