package ortb

import (
	"context"
	"net/http"
	"time"

	"clickyab.com/crane/models/campaign"

	"clickyab.com/crane/demand/layers/output/pixel"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/xlog"
)

var conversionPath = "/conversion"

// conversionHandler is the route for rtb input layer
func conversionHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	acid := r.URL.Query().Get("action_id")
	t := time.Now()
	rh := r.URL.Query().Get("imp_id")
	if rh == "" {
		_ = pixel.Render(ctx, w, nil)
		return
	}

	err := campaign.ConversionByRH(ctx, rh, acid, t, framework.RealIP(r))
	if err != nil {
		w.Header().Add("err", err.Error())
		xlog.GetWithError(ctx, err).Debug(err)
	}

	_ = pixel.Render(ctx, w, nil)
}
