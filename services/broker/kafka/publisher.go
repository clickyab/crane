package kafka

import (
	"context"
	"strings"
	"sync"
	"time"

	"clickyab.com/exchange/services/assert"
	base "clickyab.com/exchange/services/broker"
	"clickyab.com/exchange/services/config"
	"clickyab.com/exchange/services/safe"

	"github.com/Shopify/sarama"
	"github.com/Sirupsen/logrus"
)

var (
	// comma separated values
	brokerList     = config.RegisterString("services.broker.kafka.broker_list", "127.0.0.1", "kafka cluster hosts")
	flushFrequency = config.RegisterDuration("services.broker.kafka.flush_frequency", 500*time.Millisecond, "kafka flush frequency")
	sasl           = config.RegisterBoolean("services.broker.kafka.sasl_auth", false, "use sasl authentication?")
	userName       = config.RegisterString("services.broker.kafka.user", "alice", "kafka user name")
	password       = config.RegisterString("services.broker.kafka.password", "alice123@A", "kafka password")
)

type cluster struct {
	async sarama.AsyncProducer

	lock sync.RWMutex
}

func (b *cluster) Publish(j base.Job) {
	safe.GoRoutine(func() {
		b.getASync().Input() <- &sarama.ProducerMessage{
			Topic:    j.Topic(),
			Key:      sarama.StringEncoder(j.Key()),
			Metadata: j.Report(),
			Value:    j,
		}
	})
}

func (b *cluster) setASync(sa sarama.AsyncProducer) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.async = sa
}

func (b *cluster) getASync() sarama.AsyncProducer {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.async
}

func (b *cluster) errorLoop(ctx context.Context) {
	d := ctx.Done()
	for {
		select {
		case err, ok := <-b.getASync().Errors():
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
		case msg, ok := <-b.getASync().Successes():
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

	b.setASync(sa)
	go b.errorLoop(ctx)
	go b.successLoop(ctx)
	safe.GoRoutine(func() {
		<-ctx.Done()
		// this is when we need to lock the async getter
		b.lock.Lock()
		defer b.lock.Unlock()
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

	if *sasl {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = *userName
		cfg.Net.SASL.Password = *password
	}
	bl := strings.Split(*brokerList, ",")
	assert.True(len(bl) >= 1, "[CONFIGBUG] one node in kafka cluster")
	async, err := sarama.NewAsyncProducer(bl, cfg)
	assert.Nil(err)
	b.setKafka(ctx, async)
}

// NewCluster This is not a good way. but for development i need this to be done in this way.
// DEPRECATED you are not allowed to call this function
func NewCluster() interface{} {
	return &cluster{}
}
