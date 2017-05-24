package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type jsonDelivery struct {
	delivery *amqp.Delivery
}

func (jd jsonDelivery) Decode(v interface{}) error {
	return json.Unmarshal(jd.delivery.Body, v)
}

func (jd jsonDelivery) Ack(multiple bool) error {
	return jd.delivery.Ack(multiple)
}

func (jd jsonDelivery) Nack(multiple, requeue bool) error {
	return jd.delivery.Nack(multiple, requeue)
}

func (jd jsonDelivery) Reject(requeue bool) error {
	return jd.delivery.Reject(requeue)
}
