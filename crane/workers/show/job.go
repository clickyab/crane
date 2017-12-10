package show

import (
	"encoding/json"
	"net"

	"time"

	"fmt"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/sirupsen/logrus"
)

const topic = "impression"

type seats struct {
	AdID         int64   `json:"ad"`
	AdSize       int     `json:"size"`
	SlotPublicID string  `json:"slot"`
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

func (j *job) process() error {
	return fmt.Errorf("TODO: Implement me")
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
		})
	}

	return j
}
