package materialize

import (
	"bytes"
	"encoding/gob"
	"net"
	"octopus/exchange"
	"services/broker"
)

type impression struct {
	src   []byte
	topic string
	key   net.IP

	imp exchange.Impression
}

// Encode encode
func (i impression) Encode() ([]byte, error) {
		themap := []interface{}{}
		themap = append(themap, impressionToMap(i.imp))
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err := enc.Encode(themap)
		if err != nil {
			return []byte{}, err
		}
		return  buffer.Bytes(), nil
}

// Length return length
func (i impression) Length() int {
	themap := []interface{}{}
	themap = append(themap, impressionToMap(i.imp))
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(themap)
	if err != nil {
		return 0
	}
	return  len(buffer.Bytes())
}

// Topic return topic
func (i impression) Topic() string {
	return "materialize"
}

// Key return key
func (i impression) Key() string {
		return i.imp.IP().String()
}

// Report report
func (i impression) Report() func(error) {
	panic("implement me")
}

func ImpressionJob(imp exchange.Impression) broker.Job {
	return impression{
		imp: imp,
	}
}
