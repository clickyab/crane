package materialize

import (
	"encoding/json"
	"octopus/exchange"
	"services/broker"
)

type impression struct {
	data map[string]interface{}
	key  string
}

// Encode encode
func (i impression) Encode() ([]byte, error) {
	return json.Marshal(i.data)

}

// Length return length
func (i impression) Length() int {
	res, _ := i.Encode()
	return len(res)
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
	panic("implement me")
}

// ImpressionJob return a broker job
func ImpressionJob(imp exchange.Impression) broker.Job {
	return impression{
		data: impressionToMap(imp),
		key:  imp.IP().String(),
	}
}
