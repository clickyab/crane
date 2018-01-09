package banner

import (
	"context"

	"net/http"

	"fmt"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/gad/src/version"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/xlog"
)

var vs = version.GetVersion()

// Render the advertise
func Render(c context.Context, w http.ResponseWriter, ctx entity.Context) error {
	w.Header().Set("crane-version", fmt.Sprint(vs.Count))
	seats := ctx.Seats()
	assert.True(len(seats) == 1)
	s := seats[0]
	switch s.WinnerAdvertise().Type() {
	case entity.AdTypeBanner:
		return renderWebBanner(w, ctx, s)
	case entity.AdTypeDynamic:
		return renderDynamicBanner(w, ctx, s)
	case entity.AdTypeVideo:
		return renderVideoBanner(w, ctx, s)
	}

	xlog.GetWithField(c, "ad_type", s.WinnerAdvertise().Type()).Error("invalid ad type")
	return fmt.Errorf("invalid ad type ")
}
