package ortb

import (
	"context"
	"net/http"
	"time"

	"github.com/clickyab/services/xlog"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/layers/output/pixel"
	"clickyab.com/crane/workers/show"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/safe"
)

const pixelPath = "/pixel/:rh/:size/:type/:subtype/:jt"

// showPixel for ads which is not rendered by us.
func showPixel(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl, err := extractor(ctx, r)
	if err != nil {
		xlog.GetWithError(ctx, err).Error("can not extract payload in pixel route", pl)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	counter := kv.NewAEAVStore(pl.ReserveHash, clickExpire.Duration()+time.Hour).IncSubKey("I", 1)
	if counter > 1 {
		// Duplicate impression!
		pl.Suspicious = 3
	}
	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetOSUserAgent(pl.UserAgent),
		builder.SetTargetHost(r.Host),
		builder.SetIPLocation(pl.IP, nil, nil),
		builder.SetProtocolByRequest(r),
		builder.SetParent(pl.Parent, pl.Ref),
		builder.SetTID(pl.TID, pl.IP, pl.UserAgent),
		builder.SetPublisher(pl.Publisher),
		builder.SetSuspicious(pl.Suspicious),
		builder.SetFatFinger(pl.FatFinger),
	}
	b = append(b, builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.Ad, pl.Bid, time.Now().Unix(), pl.CPM, pl.SCPM, pl.requestType))
	// Build context
	c, err := builder.NewContext(b...)
	if err != nil {
		xlog.GetWithError(ctx, err).Error("context builder error in pixel rout", c)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp, _ := context.WithTimeout(ctx, 10*time.Second)
	safe.GoRoutine(exp, func() {
		job := show.NewImpressionJob(c, c.Seats()...)
		assert.Nil(broker.Publish(job))
	})

	// add capping
	safe.GoRoutine(ctx, func() {
		setCapping(ctx, pl, c.Protocol().String())
	})

	assert.Nil(pixel.Render(ctx, w, c))
}
