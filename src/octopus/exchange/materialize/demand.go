package materialize

import (
	"encoding/json"
	"octopus/exchange"
	"services/broker"
)

type demand struct {
	imp exchange.Impression
	dmn exchange.Demand
	ads map[string]exchange.Advertise

	src []byte
}

// Encode encode
func (d demand) Encode() ([]byte, error) {
	if d.src == nil {
		themap := []interface{}{}
		advertizes := []map[string]interface{}{}
		for i := range d.ads {
			advertizes = append(advertizes, advertiseToMap(d.ads[i]))
		}
		themap = append(themap, demandToMap(d.dmn), impressionToMap(d.imp, d.ads), advertizes)
		d.src, _ = json.Marshal(themap)
	}

	return d.src, nil
}

// Length return length
func (d demand) Length() int {
	x, _ := d.Encode()
	return len(x)
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
	return func(error) {}
}

// DemandJob returns a job for demand
// TODO : add a duration to this. for better view this is important
func DemandJob(imp exchange.Impression, dmn exchange.Demand, ads map[string]exchange.Advertise) broker.Job {
	return &demand{
		imp: imp,
		dmn: dmn,
		ads: ads,
	}
}
