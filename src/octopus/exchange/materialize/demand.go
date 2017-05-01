package materialize

import (
	"encoding/json"
	"octopus/exchange"
	"services/broker"
)

type demand struct {
	imp exchange.Impression
	dmn exchange.Demand
	ads []exchange.Advertise
}

// Encode encode
func (d demand) Encode() ([]byte, error) {
	themap := []interface{}{}
	advertizes := []map[string]interface{}{}
	for i := range d.ads {
		advertizes = append(advertizes, advertiseToMap(d.ads[i]))
	}
	themap = append(themap, demandToMap(d.dmn), impressionToMap(d.imp), advertizes)
	return json.Marshal(themap)
}

// Length return length
func (d demand) Length() int {
	res, _ := d.Encode()
	return len(res)
}

// Topic return topic
func (d demand) Topic() string {
	return "demand"
}

// Key return key
func (d demand) Key() string {
	return d.imp.IP().String()
}

// Report report
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
