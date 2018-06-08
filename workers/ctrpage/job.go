package ctrpage

import (
	"context"
	"encoding/json"

	"clickyab.com/crane/models/ads/statistics/locationctr"
	"clickyab.com/crane/models/pages"

	"clickyab.com/crane/models/ads"
	"clickyab.com/crane/models/seats"
	m "clickyab.com/crane/workers/models"
	"github.com/clickyab/services/xlog"
	"github.com/sirupsen/logrus"
)

const topic = "ctrpage"

//Job is an ctrpage (impression) job
//TODO: add new interface job and unexport this struct
type Job struct {
	m.Impression
	Seats []m.Seat `json:"s"`
}

// Encode this job into a byte to send over broker
func (j *Job) Encode() ([]byte, error) {
	return json.Marshal(j)
}

// Length is not required :) its here for some broker like kafka that we are not using
func (j *Job) Length() int {
	i, _ := j.Encode()
	return len(i)
}

// Topic is the job topic
func (j *Job) Topic() string {
	return topic
}

// Key is partitioning key, and not work in rabbitmq, so let it be
func (j *Job) Key() string {
	return topic
}

func (j *Job) rep(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

//Report return job report
func (j *Job) Report() func(error) {
	return j.rep
}

func (j *Job) process(ctx context.Context) error {
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
		seat, err := seats.AddAndGetSeat(j.Impression, s)
		if err != nil {
			xlog.GetWithError(ctx, err)
			errs.add(err)
			continue
		}

		pubPage, err := pages.AddAndGetPublisherPage(j.Impression)
		if err != nil {
			xlog.GetWithError(ctx, err)
			errs.add(err)
			continue
		}

		location, err := locationctr.AddAndGetCreativePerLocation(*pubPage, *seat, s.AdID, int64(s.AdSize))
		if err != nil {
			xlog.GetWithError(ctx, err)
			errs.add(err)
			continue
		}

		err = ads.UpdateImpressionLocation(j.Impression.ID, seat.ID, pubPage.ID, location.CreativeLocationID())
		if err != nil {
			xlog.GetWithError(ctx, err)
			errs.add(err)
		}
	}
	return errs.result()
}

type errorProcess struct {
	tasks  int
	errors []error
}

func (e *errorProcess) add(a ...error) {
	e.errors = append(e.errors, a...)
}

//Error return job errors
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
