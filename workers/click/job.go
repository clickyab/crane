package click

import (
	"encoding/json"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/ads"
	worker "clickyab.com/crane/workers/models"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/config"
	"github.com/sirupsen/logrus"
)

const topic = "click"

var (
	fastClick = config.RegisterInt("crane.worker.fast_click", 2, "minimum time between impression and click")
)

// job is an show (impression) job
type job struct {
	worker.Impression
	worker.Seat

	OS   entity.OS `json:"os"`
	Fast int64     `json:"fast"`
}

// Encode this job into a byte to send over broker
func (j *job) Encode() ([]byte, error) {
	return json.Marshal(j)
}

// Length is not required :) its here for some broker like kafka that we are not using
func (j *job) Length() int {
	i, _ := j.Encode()
	return len(i)
}

// Topic is the job topic
func (j *job) Topic() string {
	return topic
}

// Key is partitioning key, and not work in rabbitmq, so let it be
func (j *job) Key() string {
	return topic
}

func (j *job) rep(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

func (j *job) Report() func(error) {
	return j.rep
}

func (j *job) process() error {
	pub, err := ads.FindPublisher(j.Supplier, j.Publisher, 0, j.Type)
	if err != nil {
		return err
	}

	return ads.AdClick(pub, j.Impression, j.Seat, j.OS, j.Fast)

}

// NewClickJob return a new job for the worker
func NewClickJob(ctx entity.Context) broker.Job {
	seats := ctx.Seats()
	assert.True(len(seats) == 1)
	s := seats[0]
	fast := ctx.Timestamp().Unix() - s.ImpressionTime().Unix()
	susp := ctx.Suspicious()
	if fast < fastClick.Int64() {
		susp = 9
	}
	j := &job{
		Impression: worker.Impression{
			IP:         ctx.IP(),
			CopID:      ctx.User().ID(),
			UserAgent:  ctx.UserAgent(),
			Suspicious: susp,
			Referrer:   ctx.Referrer(),
			ParentURL:  ctx.Parent(),
			Publisher:  ctx.Publisher().Name(),
			Supplier:   ctx.Publisher().Supplier().Name(),
			Type:       ctx.SubType(),
			Timestamp:  ctx.Timestamp(),
		},
		Seat: worker.Seat{
			AdID:         s.WinnerAdvertise().ID(),
			AdSize:       s.Size(),
			SlotPublicID: s.PublicID(),
			WinnerBID:    s.Bid(),
			ReserveHash:  s.ReservedHash(),
		},
		OS:   ctx.OS(),
		Fast: fast,
	}

	return j
}
