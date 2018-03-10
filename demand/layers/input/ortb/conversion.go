package ortb

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"clickyab.com/crane/demand/layers/output/pixel"
	"clickyab.com/crane/models/ads"
	"github.com/clickyab/services/safe"
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
	// TODO : Remove the ByID function after removing the old show ad
	if impid, err := strconv.ParseInt(rh, 10, 64); err == nil {
		safe.GoRoutine(ctx, func() { ads.ConversionByID(ctx, impid, acid, t) })
	} else {
		safe.GoRoutine(ctx, func() { ads.ConversionByRH(ctx, rh, acid, t) })
	}

	_ = pixel.Render(ctx, w, nil)
}
