package ortb

import (
	"context"
	"net/http"
	"time"

	"clickyab.com/crane/demand/builder"
	"clickyab.com/crane/demand/layers/output/pixel"
	"clickyab.com/crane/workers/show"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/kv"
	"github.com/clickyab/services/safe"
	"github.com/clickyab/services/xlog"
)

const noticePath = "/notice/:rh/:size/:type/:subtype/:jt"

//pl, err := extractor(ctx, r)
//if err != nil {
//w.WriteHeader(http.StatusBadRequest)
//return
//}
//// Build context
//c, err := builder.NewContext(
//builder.SetPublisher(pl.Publisher),
//builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.AdID, pl.CpID, pl.CpAdID, pl.cpn,
//pl.Bid, pl.PreviousTime, pl.CPM, pl.SCPM, pl.requestType, pl.targetURL),
//)
//if err != nil {
//w.WriteHeader(http.StatusBadRequest)
//return
//}
//exp, _ := context.WithTimeout(ctx, 10*time.Second)
//
//safe.GoRoutine(exp, func() {
//	_ = notice.NewNoticeJob(c, c.Seats()...)
//	// TODO: uncomment below lines when notice worker is ready
//	//job := notice.NewNoticeJob(c, c.Seats()...)
//	//broker.Publish(job)
//})
//
//assert.Nil(pixel.Render(ctx, w, c))

// notice is handler for notice ad request
func noticeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	pl, err := extractor(ctx, r)
	if err != nil {
		w.Header().Add("err", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	xlog.GetWithField(ctx, "TID:", pl.TID).Debug()
	counter := kv.NewAEAVStore(pl.ReserveHash, clickExpire.Duration()+time.Hour).IncSubKey("I", 1)
	if counter > 1 {
		// Duplicate impression!
		pl.Suspicious = 3
	}
	b := []builder.ShowOptionSetter{
		builder.SetTimestamp(),
		builder.SetOSUserAgent(pl.UserAgent),
		builder.SetTargetHost(r.Host),
		builder.SetIPLocation(pl.IP, nil, nil, nil),
		builder.SetProtocolByRequest(r),
		builder.SetParent(pl.Parent, pl.Ref),
		builder.SetTID(pl.TID, pl.Did),
		builder.SetUser(nil),
		builder.SetPublisher(pl.Publisher),
		//builder.SetSuspicious(pl.Suspicious),
		builder.SetFatFinger(pl.FatFinger),
		builder.SetAdID(pl.AdID),
		builder.SetCpID(pl.CpID),
		builder.SetCpAdID(pl.CpAdID),
	}
	b = append(b, builder.SetFullSeats(pl.PublicID, pl.Size, pl.ReserveHash, pl.AdID, pl.CpID, pl.CpAdID, pl.cpn,
		pl.Bid, time.Now().Unix(), pl.CPM, pl.SCPM, pl.requestType, pl.targetURL))
	// Build context
	c, err := builder.NewContext(b...)
	if err != nil {
		w.Header().Add("err", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	exp, _ := context.WithTimeout(ctx, 10*time.Second)
	safe.GoRoutine(exp, func() {
		job := show.NewImpressionJob(c, c.Seats()...)
		broker.Publish(job)
	})

	// add capping
	safe.GoRoutine(ctx, func() {
		setCapping(ctx, pl, c.Protocol().String())
	})

	assert.Nil(pixel.Render(ctx, w, c))

}
