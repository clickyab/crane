package mocks

import "encoding/json"

type JsonDelivery struct {
	Data []byte
}

func (jd JsonDelivery) Decode(v interface{}) error {
	return json.Unmarshal(jd.Data, v)
}

func (jd JsonDelivery) Ack(multiple bool) error {
	panic("not needed")
}

func (jd JsonDelivery) Nack(multiple, requeue bool) error {
	panic("not needed")
}

func (jd JsonDelivery) Reject(requeue bool) error {
	panic("not needed")
}
