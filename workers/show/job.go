package show

import (
	"context"
	"encoding/json"
	"time"

	"clickyab.com/crane/workers/ctrpage"

	"clickyab.com/crane/models/ads/statistics/locationctr"
	"clickyab.com/crane/models/pages"

	"clickyab.com/crane/models/seats"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/models/ads"
	m "clickyab.com/crane/workers/models"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/broker"
	"github.com/clickyab/services/safe"
	"github.com/clickyab/services/xlog"
	"github.com/sirupsen/logrus"
)

const topic = "impression"

// job is an show (impression) job
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
	// TODO : multiple seat per one query
	errs := errorProcess{
		tasks: len(j.Seats),
	}
	pub, err := ads.FindPublisher(j.Supplier, j.Publisher, 0, j.PublisherType)
	if err != nil {
		return err
	}
	j.Impression.PublisherID = pub.ID()
	for _, s := range j.Seats {
		crSize := int64(s.AdSize)

		seat := seats.GetSeatByKeys(
			j.Impression.Supplier,
			s.SlotPublicID,
			j.Impression.Publisher,
			crSize,
		)
		if seat != nil {
			page := pages.GetByURLAndDomain(
				j.Impression.Publisher,
				j.Impression.ParentURL,
			)

			if page != nil {
				crlocation := locationctr.GetCRPerLocationByKeys(
					j.Impression.Publisher,
					page.ID,
					seat.ID,
					s.AdID,
					crSize,
				)

				if crlocation != nil {
					j.Impression.SeatID = seat.ID
					j.Impression.PublisherPageID = page.ID
					j.Impression.CreativesLocationID = crlocation.CreativeLocationID()
				}
			}
		}

		impID, err := ads.AddImpression(pub, j.Impression, s)

		if err != nil {
			xlog.GetWithError(ctx, err)
			errs.add(err)
		} else if j.Impression.CreativesLocationID == 0 {
			j.Impression.ID = impID

			exp, _ := context.WithTimeout(ctx, 10*time.Second)
			safe.GoRoutine(exp, func() {
				newJob := ctrpage.Job{
					Impression: j.Impression,
					Seats:      j.Seats,
				}
				broker.Publish(&newJob)
			})
		}
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
			PublisherID:   ctx.Publisher().ID(),
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
