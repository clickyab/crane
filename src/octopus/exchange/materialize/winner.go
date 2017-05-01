package materialize

import (
	"encoding/json"
	"octopus/exchange"
	"services/broker"
)

type winner struct {
	data map[string]interface{}
	key  string
}

// Encode encode
func (w *winner) Encode() ([]byte, error) {
	return json.Marshal(w.data)
}

// Length return length
func (w *winner) Length() int {
	res, _ := w.Encode()
	return len(res)
}

// Topic return topic
func (w *winner) Topic() string {
	panic("winner")
}

// Key return key
func (w *winner) Key() string {
	return w.key
}

// Report report
func (w *winner) Report() func(error) {
	panic("implement me")
}

// WinnerJob return a broker job
func WinnerJob(imp exchange.Impression, dmn exchange.Demand, ad exchange.Advertise, slotID string) broker.Job {
	return &winner{
		data: winnerToMap(imp, dmn, ad, slotID),
		key:  imp.IP().String(),
	}
}
