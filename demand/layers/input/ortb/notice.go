package ortb

import (
	"context"
	"net/http"
	"time"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/layers/output/pixel"
	"clickyab.com/crane/workers/notice"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/safe"
)

const noticePath = "/notice/:rh/:size/:type/:subtype/:jt"

// notice is handler for notice ad requestType
func noticeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pl, err := extractor(ctx, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Build context
	c, err := builder.NewContext(
		builder.SetPublisher(pl.Publisher),
		builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.Ad, pl.Bid, pl.PreviousTime, pl.CPM, pl.SCPM, pl.requestType),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp, _ := context.WithTimeout(ctx, 10*time.Second)

	safe.GoRoutine(exp, func() {
		job := notice.NewNoticeJob(c)
		broker.Publish(job)
	})

	assert.Nil(pixel.Render(ctx, w, c))

}
