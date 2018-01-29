package ortb

import (
	"context"
	"net/http"
	"time"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/layers/output/pixel"
	"clickyab.com/crane/workers/show"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
)

const pixelPath = "/pixel/:rh/:size/:type/:subtype/:jt"

// showPixel for ads which is not rendered by us.
func showPixel(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl, err := extractor(ctx, r)
	if err != nil {
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
		builder.SetIPLocation(pl.IP),
		builder.SetProtocolByRequest(r),
		builder.SetParent(pl.Parent, pl.Ref),
		builder.SetTID(pl.TID, pl.IP, pl.UserAgent),
		builder.SetType(pl.Type, pl.SubType),
		builder.SetPublisher(pl.Publisher),
		builder.SetSuspicious(pl.Suspicious),
		builder.SetFatFinger(pl.FatFinger),
	}
	if pl.Type == entity.RequestTypeDemand {
		b = append(b, builder.DoNotShowTFrame())
	}
	b = append(b, builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.Ad, pl.Bid, time.Now().Unix(), pl.CPM, pl.SCPM))
	// Build context
	c, err := builder.NewContext(b...)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp, _ := context.WithTimeout(ctx, 10*time.Second)
	safe.GoRoutine(exp, func() {
		job := show.NewImpressionJob(c, c.Seats()...)
		broker.Publish(job)
	})

	assert.Nil(pixel.Render(ctx, w, c))
}
