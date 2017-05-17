package jsonbackend

import (
	"encoding/json"
	"octopus/exchange"
	"services/broker"

	"github.com/Sirupsen/logrus"
)

type impression struct {
	data map[string]interface{}
	key  string

	src []byte
}

// Encode encode
func (i impression) Encode() ([]byte, error) {
	if i.src == nil {
		i.src, _ = json.Marshal(i.data)
	}

	return i.src, nil

}

// Length return length
func (i impression) Length() int {
	x, _ := i.Encode()
	return len(x)
}

// Topic return topic
func (i impression) Topic() string {
	return "impression"
}

// Key return key
func (i impression) Key() string {
	return i.key
}

// Report report
func (i impression) Report() func(error) {
	return func(err error) {
		if err != nil {
			logrus.Warn(err)
		}
	}
}

// ImpressionJob return a broker job
func ImpressionJob(imp exchange.Impression) broker.Job {
	return impression{
		data: impressionToMap(imp, nil),
		key:  imp.IP().String(),
	}
}
