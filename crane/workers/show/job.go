package show

import (
	"encoding/json"
	"net"

	"time"

	"fmt"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/broker"
	"github.com/sirupsen/logrus"
)

const topic = "impression"

// job is an show (impression) job
type job struct {
	IP           net.IP             `json:"ip"`
	AdID         int64              `json:"ad"`
	AdSize       int                `json:"size"`
	CopID        string             `json:"cop"`
	SlotPublicID string             `json:"slot"`
	UserAgent    string             `json:"ua"`
	WinnerBID    float64            `json:"wb"`
	Suspicious   bool               `json:"sp"`
	Referrer     string             `json:"r"`
	ParentURL    string             `json:"par"`
	Publisher    int64              `json:"pub"`
	Type         entity.RequestType `json:"t"`

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
func NewImpressionJob(ctx entity.Context, seat entity.Seat) broker.Job {
	j := &job{
		IP:           ctx.IP(),
		AdID:         seat.WinnerAdvertise().ID(),
		AdSize:       seat.Size(),
		CopID:        ctx.User().ID(),
		SlotPublicID: seat.PublicID(),
		UserAgent:    ctx.UserAgent(),
		WinnerBID:    seat.Bid(),
		Suspicious:   false, // TODO : After adding to system
		Referrer:     ctx.Referrer(),
		ParentURL:    ctx.Parent(),
		Publisher:    ctx.Publisher().ID(),
		Type:         ctx.Type(),
		Timestamp:    ctx.Timestamp(),
	}

	return j
}
