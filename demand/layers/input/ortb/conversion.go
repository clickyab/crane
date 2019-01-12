package ortb

import (
	"context"
	"net/http"
	"time"

	"clickyab.com/crane/demand/layers/output/pixel"
	"clickyab.com/crane/models/ads"
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

	err := ads.ConversionByRH(ctx, rh, acid, t)
	if err == nil {
		xlog.GetWithError(ctx, err).Debug(err)

	}

	_ = pixel.Render(ctx, w, nil)
}
