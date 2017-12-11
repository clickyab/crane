package show

import (
	"context"
	"encoding/json"
	"net"

	"time"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/models"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/xlog"
	"github.com/sirupsen/logrus"
)

const topic = "impression"

type seats struct {
	AdID         int64   `json:"ad"`
	AdSize       int     `json:"size"`
	SlotPublicID string  `json:"slot"`
	ReserveHash  string  `json:"reserve_hash"`
	WinnerBID    float64 `json:"wb"`
}

// job is an show (impression) job
type job struct {
	IP         net.IP             `json:"ip"`
	CopID      string             `json:"cop"`
	UserAgent  string             `json:"ua"`
	Suspicious int                `json:"sp"`
	Referrer   string             `json:"r"`
	ParentURL  string             `json:"par"`
	Publisher  string             `json:"pub"`
	Supplier   string             `json:"sub"`
	Type       entity.RequestType `json:"t"`
	Alexa      bool               `json:"a"`

	Seats []seats `json:"s"`

	Timestamp time.Time `json:"ts"`
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

func (j *job) process(ctx context.Context) error {
	errs := errorProcess{
		tasks: len(j.Seats),
	}
	for _, v := range j.Seats {
		err := models.AddImpression(j.Supplier, v.ReserveHash, j.Publisher, j.Referrer, j.ParentURL,
			v.SlotPublicID, j.CopID, v.AdSize, j.Suspicious, v.AdID, j.IP, v.WinnerBID, j.Alexa, j.Timestamp, j.Type)
		if err != nil {
			xlog.GetWithError(ctx, err)
			errs.add(err)
		}
	}
	return errs.result()
}

// NewImpressionJob return a new job for the worker
func NewImpressionJob(ctx entity.Context, s ...entity.Seat) broker.Job {
	assert.True(len(s) > 0)
	j := &job{

		IP:         ctx.IP(),
		CopID:      ctx.User().ID(),
		UserAgent:  ctx.UserAgent(),
		Suspicious: ctx.Suspicious(),
		Referrer:   ctx.Referrer(),
		ParentURL:  ctx.Parent(),
		Publisher:  ctx.Publisher().Name(),
		Supplier:   ctx.Publisher().Supplier().Name(),
		Type:       ctx.Type(),
		Timestamp:  ctx.Timestamp(),
		Alexa:      ctx.Alexa(),
	}
	for i := range s {
		j.Seats = append(j.Seats, seats{
			AdID:         s[i].WinnerAdvertise().ID(),
			AdSize:       s[i].Size(),
			SlotPublicID: s[i].PublicID(),
			WinnerBID:    s[i].Bid(),
			ReserveHash:  s[i].ReservedHash(),
		})
	}

	return j
}

type errorProcess struct {
	tasks  int
	errors []error
}

func (e *errorProcess) add(a ...error) {
	e.errors = append(e.errors, a...)
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
