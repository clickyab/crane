package materialize

import (
	"bytes"
	"encoding/gob"
	"net"
	"octopus/exchange"
	"services/broker"
)

type demand struct {
	key    net.IP
	src    []byte
	topic  string
	lenght int

	imp exchange.Impression
	dmn exchange.Demand
	ads []exchange.Advertise
}

func (d *demand) Encode() ([]byte, error) {
	if d.src == nil {
		themap := []interface{}{}
		advertizes := []map[string]interface{}{}
		for i := range d.ads {
			advertizes = append(advertizes, advertiseToMap(d.ads[i]))
		}
		themap = append(themap, demandToMap(d.dmn), impressionToMap(d.imp), advertizes)

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(themap)
		if err != nil {
			return []byte{}, err
		}

		d.src = buf.Bytes()
	}
	return d.src, nil
}

func (d *demand) Length() int {
	if len(d.src) == 0 {
		var err error
		d.src, err = d.Encode()
		if err != nil {
			panic("asd")
		}
	}
	return len(d.src)
}

func (d *demand) Topic() string {
	return "materialize"
}

func (d *demand) Key() string {
	if d.key == nil {
		d.key = d.imp.IP()
	}
	return d.key.String()
}

func (d *demand) Report() func(error) {
	panic("implement me")
}

// DemandJob returns a job for demand
func DemandJob(imp exchange.Impression, dmn exchange.Demand, ads []exchange.Advertise) broker.Job {
	return demand{
		imp: imp,
		dmn: dmn,
		ads: ads,
	}
}
