package restful

import (
	"assert"
	"bytes"
	"context"
	"encoding/json"
	"entity"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"demands/models"

	"github.com/Sirupsen/logrus"
)

type demand struct {
	client *http.Client

	endPoint           string
	winPoint           *url.URL
	maxIdleConnections int

	requestTimeout time.Duration
	key            string

	minuteLimit int64
	hourLimit   int64
	dayLimit    int64
	weekLimit   int64
	monthLimit  int64
}

func (d *demand) Name() string {
	return d.key
}

func (d *demand) checkLimits() bool {

}

func (d *demand) Provide(ctx context.Context, imp entity.Impression, ch chan map[string]entity.Advertise) {
	defer close(ch)
	tmp := impressionToMap(imp)
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(tmp); err != nil {
		logrus.Debug(err)
		return
	}
	req, err := http.NewRequest("POST", d.endPoint, bytes.NewBuffer(buf.Bytes()))
	if err != nil {
		logrus.Debug(err)
		return
	}
	resp, err := d.client.Do(req.WithContext(ctx))
	if err != nil {
		logrus.Debug(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Debugf("status code is %d", resp.StatusCode)
		return
	}

	ads := map[string]restAd{}
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err := dec.Decode(&tmp); err != nil {
		logrus.Debug(err)
		return
	}

	adsInter := make(map[string]entity.Advertise, len(ads))
	for i := range ads {
		tmp := ads[i]
		// set demand for win resp
		tmp.demand = d
		adsInter[i] = tmp
	}

	ch <- adsInter
}

func (d *demand) Status(c context.Context, h http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (d *demand) Win(ctx context.Context, id string, cpm int64) {
	u := *d.winPoint
	u.Query().Add("win", id)
	u.Query().Add("cpm", fmt.Sprint(cpm))
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		logrus.Debug(err)
		return
	}
	resp, err := d.client.Do(req.WithContext(ctx))
	if err != nil {
		logrus.Debug(err)
		return
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusNoContent {
		return
	}

	logrus.Debug("winner call status was %d", resp.StatusCode)

}

func (d *demand) createConnection() {
	d.client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: d.maxIdleConnections,
		},
		Timeout: d.requestTimeout,
	}
}

// NewRestfulClient return a new client for restful call
func NewRestfulClient(d models.Demand) entity.Demand {
	win, err := url.Parse(d.WinURL)
	assert.Nil(err)
	return &demand{
		endPoint:           d.GetURL,
		winPoint:           win,
		maxIdleConnections: d.IdleConnections,
		requestTimeout:     d.GetTimeout(),
		key:                d.Name,
		minuteLimit:        d.MinuteLimit,
		hourLimit:          d.HourLimit,
		dayLimit:           d.DayLimit,
		weekLimit:          d.WeekLimit,
		monthLimit:         d.MonthLimit,
	}
}
