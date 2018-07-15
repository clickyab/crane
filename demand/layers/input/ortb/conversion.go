package ortb

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"clickyab.com/crane/demand/layers/output/pixel"
	"clickyab.com/crane/models/ads"
)

var conversionPath = "/conversion"

// conversionHandler is the route for rtb input layer
func conversionHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	acid := r.URL.Query().Get("action_id")
	var found bool
	t := time.Now()

	rh := r.URL.Query().Get("imp_id")
	if rh == "" {
		_ = pixel.Render(ctx, w, nil, found)
		return
	}
	// TODO : Remove the ByID function after removing the old show ad
	if impid, err := strconv.ParseInt(rh, 10, 64); err == nil {
		err = ads.ConversionByID(ctx, impid, acid, t)
		if err == nil {
			found = true
		}
	} else {
		err = ads.ConversionByRH(ctx, rh, acid, t)
		if err == nil {
			found = true
		}
	}

	_ = pixel.Render(ctx, w, nil, found)
}
