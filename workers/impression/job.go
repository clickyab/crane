package impression

import (
	"context"
	"encoding/json"
	"time"

	"github.com/clickyab/services/config"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/ads"
	m "clickyab.com/crane/workers/models"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/sirupsen/logrus"
)

var bulkCount = config.RegisterInt64("crane.workers.impressions.bulk.count", 200, "expire time of ads")
var bulkTime = config.RegisterDuration("crane.workers.impressions.bulk.time", time.Second*2, "expire time of ads")

// job is an impression (impression) job
type job struct {
	m.Impression
	Seats []m.Seat `json:"s"`
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
	return "impression"
}

// Key is partitioning key, and not work in rabbitmq, so let it be
func (j *job) Key() string {
	return "impression"
}

func (j *job) rep(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

func (j *job) Report() func(error) {
	return j.rep
}

var impressions = make(chan m.Impression)

func (j *job) process(ctx context.Context) error {
	// TODO : multiple seat per one query
	errs := errorProcess{
		tasks: len(j.Seats),
	}
	pub, err := ads.FindPublisher(j.Supplier, j.Publisher, 0, j.PublisherType)
	if err != nil {
		return err
	}
	for i := range j.Seats {
		j.Impression.Pub = pub
		j.Impression.Seat = j.Seats[i]
		impressions <- j.Impression
	}
	return errs.result()
}

// NewImpressionJob return a new job for the worker
func NewImpressionJob(ctx entity.Context, s ...entity.Seat) broker.Job {
	assert.True(len(s) > 0)
	j := &job{
		Impression: m.Impression{
			IP:            ctx.IP(),
			CopID:         ctx.User().ID(),
			UserAgent:     ctx.UserAgent(),
			Suspicious:    ctx.Suspicious(),
			Referrer:      ctx.Referrer(),
			ParentURL:     ctx.Parent(),
			Publisher:     ctx.Publisher().Name(),
			Supplier:      ctx.Publisher().Supplier().Name(),
			Timestamp:     ctx.Timestamp(),
			PublisherType: ctx.Publisher().Type(),
		},
	}
	for i := range s {
		j.Seats = append(j.Seats, m.Seat{
			AdID:         s[i].WinnerAdvertise().ID(),
			AdSize:       s[i].Size(),
			SlotPublicID: s[i].PublicID(),
			WinnerBID:    s[i].Bid(),
			ReserveHash:  s[i].ReservedHash(),
			CPM:          s[i].CPM(),
			SCPM:         s[i].SupplierCPM(),
			Type:         s[i].RequestType(),
		})
	}

	return j
}

type errorProcess struct {
	tasks  int
	errors []error
}

func (e *errorProcess) Error() string {
	res := ""
	for i := range e.errors {
		res += e.errors[i].Error() + "\n"
	}
	return res
}
func (e *errorProcess) result() error {
	if len(e.errors) >= e.tasks {
		return e
	}
	return nil
}
