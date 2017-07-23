package native

import (
	"testing"

	"regexp"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/entity/mock_entity"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/goconvey/convey"
)

var regex = regexp.MustCompile(`<div class="cyb-suggest`)

func TestRenderer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	convey.Convey("native output render check", t, func() {
		mockSlots := make([]entity.Slot, 0)
		for _, i := range impMockDetails.Slots {
			mockAd := mock_entity.NewMockAdvertise(ctrl)
			mockAd.EXPECT().Type().Return(i.AdType).AnyTimes()
			mockAd.EXPECT().Attributes().Return(i.AdAttr).AnyTimes()
			mockAd.EXPECT().TargetURL().Return(i.AdTargetURL).AnyTimes()
			mockAd.EXPECT().Media().Return(i.AdMedia).AnyTimes()
			mockSlot := mock_entity.NewMockSlot(ctrl)
			mockSlot.EXPECT().WinnerAdvertise().Return(mockAd).AnyTimes()
			mockSlot.EXPECT().ID().Return(i.AdID).AnyTimes()

			mockSlots = append(mockSlots, mockSlot)
		}

		imp := mock_entity.NewMockImpression(ctrl)
		imp.EXPECT().Slots().Return(mockSlots).AnyTimes()
		imp.EXPECT().Attributes().Return(impMockDetails.ImpAttr).AnyTimes()

		writer := &m{}
		NewRenderer().Render(writer, imp, clickProvider{})

		convey.So(len(regex.FindAll([]byte(writer.data), -1)), should.Equal, len(impMockDetails.Slots)+1)
	})

}

type m struct {
	data string
}

func (m *m) Write(p []byte) (n int, err error) {
	m.data = string(p)
	return 0, nil
}

type clickProvider struct{}

func (clickProvider) ClickURL(s entity.Slot, i entity.Impression) string {
	return ""
}

type impEntity struct {
	ImpAttr map[string]string
	Slots   []slotEntity
}

type slotEntity struct {
	AdType      entity.AdType
	AdAttr      map[string]interface{}
	AdID        string
	AdTargetURL string
	AdMedia     string
}

var impMockDetails = impEntity{
	ImpAttr: map[string]string{
		"title":    "imp title",
		"style":    "imp style",
		"fontSize": "imp fontSize",
		"position": "imp position",
	},
	Slots: []slotEntity{
		{
			AdType: entity.AdTypeNative,
			AdAttr: map[string]interface{}{
				"title": "ad1title",
			},
			AdID:        "ad1ID",
			AdTargetURL: "ad1TargetURL",
			AdMedia:     "ad1Source",
		},

		{
			AdType: entity.AdTypeNative,
			AdAttr: map[string]interface{}{
				"title": "ad2title",
			},
			AdID:        "ad2ID",
			AdTargetURL: "ad2TargetURL",
			AdMedia:     "ad2Source",
		},

		{
			AdType: entity.AdTypeNative,
			AdAttr: map[string]interface{}{
				"title": "ad3title",
			},
			AdID:        "ad3ID",
			AdTargetURL: "ad3TargetURL",
			AdMedia:     "ad3Source",
		},
	},
}
