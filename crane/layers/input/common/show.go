package common

import (
	"context"
	"net/http"
	"time"

	"clickyab.com/crane/crane/builder"
	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/output/banner"
	"clickyab.com/crane/crane/workers/show"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
)

const showPath = "/banner/:rh/:size/:type/:jt"

// show is handler for show ad request
func showBanner(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl, err := extractor(ctx, r)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetOSUserAgent(pl.UserAgent),
		builder.SetRequest(r.Host, r.Method),
		builder.SetIPLocation(pl.IP),
		builder.SetAlexa(pl.UserAgent, http.Header{}),
		builder.SetProtocolByRequest(r),
		builder.SetParent(pl.Parent, pl.Ref),
		builder.SetTID(pl.TID, pl.IP, pl.UserAgent),
		builder.SetType(pl.Type),
		builder.SetPublisher(pl.Publisher),
		builder.SetSuspicious(pl.Suspicious),
	}
	if pl.Type == entity.RequestTypeDemand {
		b = append(b, builder.DoNotShowTFrame())
	}
	b = append(b, builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.Ad, pl.Bid))
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
