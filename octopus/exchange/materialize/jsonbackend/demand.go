package jsonbackend

import (
	"encoding/json"

	"clickyab.com/exchange/octopus/exchange"
	"github.com/clickyab/services/broker"

	"github.com/Sirupsen/logrus"
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
		themap := make(map[string]interface{})
		themap["demand"] = demandToMap(d.dmn)
		themap["impression"] = impressionToMap(d.imp, d.ads)
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
	return func(err error) {
		if err != nil {
			logrus.Warn(err)
		}
	}
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
