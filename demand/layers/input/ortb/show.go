package ortb

import (
	"context"
	"net/http"
	"time"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/demand/layers/output/banner"
	"clickyab.com/crane/workers/show"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
)

const showPath = "/banner/:rh/:size/:type/:subtype/:jt"

// show is handler for show ad requestType
func showBanner(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl, err := extractor(ctx, r)
	if err != nil {
		logrus.Error(err)
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
		builder.SetNoTiny(!pl.Tiny),
		builder.SetTargetHost(r.Host),
		builder.SetIPLocation(pl.IP, nil, nil),
		builder.SetProtocolByRequest(r),
		builder.SetParent(pl.Parent, pl.Ref),
		builder.SetTID(pl.TID, pl.IP, pl.UserAgent),
		builder.SetPublisher(pl.Publisher),
		builder.SetSuspicious(pl.Suspicious),
		builder.SetFatFinger(pl.FatFinger),
	}
	if pl.Type == entity.InputTypeDemand {
		b = append(b, builder.DoNotShowTFrame())
	}
	b = append(b, builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.Ad, pl.Bid, time.Now().Unix(), pl.CPM, pl.SCPM, pl.requestType))
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

	assert.Nil(banner.Render(ctx, w, c))
}
