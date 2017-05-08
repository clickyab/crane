package restful

import (
	"context"
	"encoding/json"
	"net/http"
	"octopus/core"
	"octopus/rtb"
	"octopus/supliers"

	"fmt"
	"services/assert"
	"services/eav"
	"time"

	"services/httplib"

	"octopus/exchange/materialize"
	"services/broker"
	"services/safe"

	"github.com/fzerorubigd/xmux"
)

func getAd(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	dec := json.NewEncoder(w)
	key := xmux.Param(ctx, "key")
	imp, err := supliers.GetImpression(key, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		dec.Encode(struct {
			Error error
		}{
			Error: err,
		})
		return
	}
	// OK push it to broker
	safe.GoRoutine(func() {
		jImp := materialize.ImpressionJob(imp)
		broker.Publish(jImp)
	}, r)
	nCtx, cnl := context.WithCancel(ctx)
	defer cnl()
	ads := core.Call(nCtx, imp)
	ads = rtb.Moderate(imp.Source(), ads)
	res := rtb.SelectCPM(imp, ads)
	// save winner in store
	for i := range res {
		store := eav.NewEavStore(res[i].TrackID())
		store.SetSubKey("IP",
			httplib.RealIP(r),
		).SetSubKey("SLOT",
			res[i].ID(),
		).SetSubKey("DEMAND",
			res[i].Demand().Name(),
		).SetSubKey("BID",
			fmt.Sprintf("%d", res[i].WinnerCPM()),
		).SetSubKey("TRACK",
			res[i].TrackID(),
		).SetSubKey("ADID",
			res[i].ID(),
		).SetSubKey("TIME", fmt.Sprint(imp.Time().Unix()))
		assert.Nil(store.Save(24 * time.Hour))
	}
	err = imp.Source().Supplier().Renderer().Render(imp, res, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		dec.Encode(struct {
			Error string
		}{
			Error: err.Error(),
		})
	}
}
