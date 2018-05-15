package rabbitmq

import (
	"fmt"
	"time"

	"github.com/clickyab/services/broker"
)

// FakeAmqp struct
type FakeAmqp struct {
}

var (
	fakeConnectionPool = make(map[string]string, 0)
	fakePublisherPool  = make(map[string]string, 0)
	jobs               = make(map[string][][]byte, 0)
)

// MakeConnections Dial to amqp and make a connection pool
func (rb *FakeAmqp) MakeConnections(count int) ([]string, error) {
	var keys []string

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("cn_%d", time.Now().UnixNano())
		keys = append(keys, key)
		fakeConnectionPool[key] = fmt.Sprintf("Fake_Connection_%s", key)
	}

	return keys, nil
}

// ExchangeDeclare Declare new exchange. to be sure exchange is declared
func (rb *FakeAmqp) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool) error {
	fmt.Printf("ok fake exchange with name %s declared", name)

	return nil
}

// MakePublisher register new job publisher
func (rb *FakeAmqp) MakePublisher(connectionKey string) (string, error) {
	key := fmt.Sprintf("pub_%d", time.Now().UnixNano())
	fakePublisherPool[key] = fmt.Sprintf("FAKE_PUBLISHER_%s", key)

	return key, nil
}

// Publish job to amqp
func (rb *FakeAmqp) Publish(in broker.Job, pubKey string) error {
	data, err := in.Encode()
	if err != nil {
		return err
	}

	jobs[in.Topic()] = append(jobs[in.Topic()], data)
	fmt.Printf("ok we publish a fake job. at topic %s we have %d job now", in.Topic(), len(jobs[in.Topic()]))

	return nil
}

// FinalizeWait is a function to wait for all publication to finish. after calling this,
// must not call the PublishEvent
func (rb *FakeAmqp) FinalizeWait() {
	return
}

// RegisterConsumer register new cunsumer and fetch some of jobs
func (rb *FakeAmqp) RegisterConsumer(consumer broker.Consumer, prefetchCount int) error {
	fmt.Printf("start consumer on job at topic %s we have %d job now", consumer.Topic(), len(jobs[consumer.Topic()]))

	jobs[consumer.Topic()] = append(jobs[consumer.Topic()][:prefetchCount], jobs[consumer.Topic()][prefetchCount+1:]...)

	fmt.Printf("ok we done %d job now we have %d job in queue", prefetchCount, len(jobs[consumer.Topic()]))
	return nil
}
