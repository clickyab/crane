package banner

import (
	"io"

	"context"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/xlog"
)

// Render the advertise
func Render(c context.Context, w io.Writer, ctx entity.Context, s entity.Seat) error {
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
