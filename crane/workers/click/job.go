package click

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

const topic = "click"

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

	AdID         int64   `json:"ad"`
	AdSize       int     `json:"size"`
	SlotPublicID string  `json:"slot"`
	WinnerBID    float64 `json:"wb"`

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

// NewClickJob return a new job for the worker
func NewClickJob(ctx entity.Context) broker.Job {
	seats := ctx.Seats()
	assert.True(len(seats) == 1)
	s := seats[0]
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

		AdID:         s.WinnerAdvertise().ID(),
		AdSize:       s.Size(),
		SlotPublicID: s.PublicID(),
		WinnerBID:    s.Bid(),
	}

	return j
}
