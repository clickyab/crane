package kafka

import (
	"testing"

	"time"

	"context"

	"errors"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

type job struct {
	data   string
	topic  string
	err    error
	called int

	done chan struct{}
}

func (j *job) Encode() ([]byte, error) {
	return []byte(j.data), nil
}

func (j *job) Length() int {
	return len(j.data)
}

func (j *job) Topic() string {
	return j.topic
}

func (j *job) Key() string {
	return "SOMEKEY"
}

func (j *job) Report() func(error) {
	return func(e error) {
		j.err = e
		j.called++
		j.done <- struct{}{}
	}
}

func TestSpec(t *testing.T) {
	b := &cluster{}
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	Convey("Test normal Publish ", t, func() {
		ctx, cl := context.WithCancel(context.Background())
		defer cl()
		async := mocks.NewAsyncProducer(t, cfg)
		async.ExpectInputAndSucceed()
		b.setKafka(ctx, async)
		d := make(chan struct{})
		j := &job{
			data:  "DATA",
			topic: "TOPIC",
			done:  d,
		}
		b.Publish(j)

		select {
		case <-time.After(10 * time.Second):
			So(nil, ShouldEqual, "THE TIME HAS PASSED.")
		case <-d:
		}

		So(j.called, ShouldEqual, 1)
		So(j.err, ShouldBeNil)
	})

	Convey("Test failed Publish ", t, func() {
		ctx, cl := context.WithCancel(context.Background())
		defer cl()

		errExpected := errors.New("this is error")
		async := mocks.NewAsyncProducer(t, cfg)
		async.ExpectInputAndFail(errExpected)
		b.setKafka(ctx, async)
		d := make(chan struct{})
		j := &job{
			data:  "DATA",
			topic: "TOPIC",
			done:  d,
		}
		b.Publish(j)

		select {
		case <-time.After(10 * time.Second):
			So(nil, ShouldEqual, "THE TIME HAS PASSED.")
		case <-d:
		}

		So(j.called, ShouldEqual, 1)
		So(j.err.Error(), ShouldEndWith, errExpected.Error())
	})
}
