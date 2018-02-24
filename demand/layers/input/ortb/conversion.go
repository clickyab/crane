package ortb

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"clickyab.com/crane/demand/layers/output/pixel"
	"clickyab.com/crane/models/ads"
	"github.com/clickyab/services/safe"
	"github.com/clickyab/services/xlog"
	"github.com/sirupsen/logrus"
)

var conversionPath = "/conversion"

func conversionByID(ctx context.Context, impid int64, acid string, t time.Time) {
	var err error
	imp, err := ads.FindImpressionByID(impid, t)
	if err == nil {
		addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
		return
	}
	imp, err = ads.FindImpressionByID(impid, t.AddDate(0, 0, -1))
	if err == nil {
		addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
		return
	}
	imp, err = ads.FindImpFromClickByImpID(impid)
	if err == nil {
		addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
		return
	}

	xlog.GetWithError(ctx, err).WithFields(logrus.Fields{
		"action_id": acid,
		"imp_id":    fmt.Sprint(impid),
	})
}

// conversionHandler is the route for rtb input layer
func conversionHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	acid := r.URL.Query().Get("action_id")
	t := time.Now()

	if impid, err := strconv.ParseInt(r.URL.Query().Get("imp_id"), 10, 64); err == nil {
		safe.GoRoutine(ctx, func() { conversionByID(ctx, impid, acid, t) })
	} else {
		rh := r.URL.Query().Get("imp_id")
		safe.GoRoutine(ctx, func() { conversionByRH(ctx, rh, acid, t) })
	}

	_ = pixel.Render(ctx, w, nil)
}

func conversionByRH(ctx context.Context, rh string, acid string, t time.Time) {
	imp, err := ads.FindImpressionByRH(rh, t)
	if err == nil {
		addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
		return
	}
	imp, err = ads.FindImpressionByRH(rh, t.AddDate(0, 0, -1))
	if err == nil {
		addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
		return
	}
	imp, err = ads.FindImpFromClickByRH(rh)
	if err == nil {
		addConversion(imp.WID.Int64, imp.AppID.Int64, imp.WpID.Int64, imp.CaID.Int64, imp.AdID.Int64, imp.CopID.Int64, imp.CpID.Int64, imp.SlotID.Int64, imp.ImpID.Int64, acid)
		return
	}

	xlog.GetWithError(ctx, err).WithFields(logrus.Fields{
		"action_id": acid,
		"imp_id":    rh,
	})

}
func addConversion(wID, appID, wpID, caID, adID, copID, cpID, slotID, impID int64, acid string) {
	_ = ads.AddConversion(wID, appID, wpID, caID, adID, copID, cpID, slotID, impID, acid)
}
