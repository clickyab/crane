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

func (d demand) Encode() ([]byte, error) {
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

		return buf.Bytes(),nil

}

func (d demand) Length() int {
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
		return 0
	}
	return len(buf.Bytes())
}

func (d demand) Topic() string {
	return "materialize"
}

func (d demand) Key() string {
	return d.imp.IP().String()
}

func (d demand) Report() func(error) {
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
