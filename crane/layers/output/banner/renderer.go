package banner

import (
	"net/http"

	"context"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/xlog"
)

type renderer struct {
}

func (renderer) Render(c context.Context, w http.ResponseWriter, ctx *builder.Context, s entity.Slot, ad entity.Advertise) error {
	switch ad.Type() {
	case entity.AdTypeBanner:
		return renderWebBanner(w, ctx, s, ad)
	case entity.AdTypeDynamic:
		return renderDynamicBanner(w, ctx, s, ad)
	case entity.AdTypeVideo:
		return renderVideoBanner(w, ctx, s, ad)
	}

	xlog.GetWithField(c, "ad_type", ad.Type()).Panic("invalid ad type")
}
