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

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/xmux"
)

// GetAd is route to get the ad from a restful endpoint
func GetAd(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
	jImp := materialize.ImpressionJob(imp)
	broker.Publish(jImp)
	nCtx, cnl := context.WithCancel(ctx)
	defer cnl()
	ads := core.Call(nCtx, imp)
	logrus.Debugf("%d ads is passed the system from exchange calls", len(ads))
	res := rtb.SelectCPM(imp, ads)
	logrus.Debugf("%d ads is passed the system select", len(res))
	// Publish them into message broker
	for i := range res {
		// TODO : why this happen?
		if res[i] != nil {
			broker.Publish(materialize.WinnerJob(
				imp,
				res[i],
				i,
			))

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
			).SetSubKey("TIME",
				fmt.Sprint(imp.Time().Unix()),
			).SetSubKey("PUBLISHER",
				imp.Source().Name(),
			).SetSubKey("SUPPLIER",
				imp.Source().Supplier().Name(),
			)
			assert.Nil(store.Save(24 * time.Hour))
		}
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
