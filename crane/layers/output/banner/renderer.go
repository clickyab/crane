package banner

import (
	"context"

	"net/http"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/xlog"
)

// Render the advertise
func Render(c context.Context, w http.ResponseWriter, ctx entity.Context) error {
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

	xlog.GetWithField(c, "ad_type", s.WinnerAdvertise().Type()).Panic("invalid ad type")
	return nil
}
