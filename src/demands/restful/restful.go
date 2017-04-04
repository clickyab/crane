package restful

import (
	"context"
	"entity"
	"net/http"
	"time"

	"bytes"
	"encoding/json"

	"net/url"

	"assert"

	"github.com/Sirupsen/logrus"
)

type demand struct {
	client *http.Client

	endPoint           string
	winPoint           *url.URL
	maxIdleConnections int

	requestTimeout time.Duration
	key            string
}

func (d *demand) Name() string {
	return d.key
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

	ads := map[string]restAd{}
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err := dec.Decode(&tmp); err != nil {
		logrus.Debug(err)
		return
	}

	adsInter := make(map[string]entity.Advertise, len(ads))
	for i := range ads {
		adsInter[i] = ads[i]
	}

	ch <- adsInter
}

func (d *demand) Status(context.Context, http.ResponseWriter, *http.Request) {
	panic("implement me")
}

func (d *demand) Win(ctx context.Context, id string) {
	u := *d.winPoint
	u.Query().Add("win", id)
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
func NewRestfulClient(name, endpoint, winpoint string, maxIdleConnection int, timeout time.Duration) entity.Demand {
	win, err := url.Parse(winpoint)
	assert.Nil(err)
	return &demand{
		endPoint:           endpoint,
		winPoint:           win,
		maxIdleConnections: maxIdleConnection,
		requestTimeout:     timeout,
		key:                name,
	}
}
