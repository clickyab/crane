package kafka

import (
	"context"
	"services/assert"
	base "services/broker"
	"services/config"
	"services/safe"
	"strings"
	"time"

	"services/initializer"

	"github.com/Shopify/sarama"
	"github.com/Sirupsen/logrus"
)

var (
	// comma separated values
	brokerList     = config.RegisterString("services.cluster.kafka.broker_list", "127.0.0.1")
	flushFrequency = config.RegisterDuration("services.cluster.kafka.flush_frequency", 500*time.Millisecond)
)

type cluster struct {
	async sarama.AsyncProducer
}

func (b *cluster) Publish(j base.Job) {
	b.async.Input() <- &sarama.ProducerMessage{
		Topic:    j.Topic(),
		Key:      sarama.StringEncoder(j.Key()),
		Metadata: j.Report(),
		Value:    j,
	}
}

func (b *cluster) errorLoop(ctx context.Context) {
	d := ctx.Done()
	for {
		select {
		case err, ok := <-b.async.Errors():
			if !ok {
				return
			}
			meta, ok := err.Msg.Metadata.(func(error))
			if !ok {
				// this is not us? why?
				logrus.Error(err)
				continue
			}
			safe.Routine(func() { meta(err) })
		case <-d:
			return
		}
	}
}

func (b *cluster) successLoop(ctx context.Context) {
	d := ctx.Done()
	for {
		select {
		case msg, ok := <-b.async.Successes():
			if !ok {
				return
			}
			meta, ok := msg.Metadata.(func(error))
			if ok && meta != nil {
				safe.Routine(func() { meta(nil) })
			}
		case <-d:
			return
		}
	}
}

func (b *cluster) setKafka(ctx context.Context, sa sarama.AsyncProducer) {
	done := ctx.Done()
	assert.NotNil(done, "[BUG] context is not cancelable")

	b.async = sa
	go b.errorLoop(ctx)
	go b.successLoop(ctx)
	safe.GoRoutine(func() {
		<-ctx.Done()
		assert.Nil(b.async.Close())
	})

}

func (b *cluster) Initialize(ctx context.Context) {

	if *flushFrequency < time.Millisecond {
		*flushFrequency = 500 * time.Millisecond
	}
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForLocal     // Only wait for the leader to ack
	cfg.Producer.Compression = sarama.CompressionSnappy // Compress messages
	cfg.Producer.Flush.Frequency = *flushFrequency      // Flush batches every 500ms
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	bl := strings.Split(*brokerList, ",")
	assert.True(len(bl) >= 1, "[CONFIGBUG] one node in kafka cluster")
	async, err := sarama.NewAsyncProducer(bl, cfg)
	assert.Nil(err)
	b.setKafka(ctx, async)
}

func init() {
	b := &cluster{}
	initializer.Register(b)
	base.SetActiveBroker(b)
}
