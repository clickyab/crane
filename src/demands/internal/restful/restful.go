package restful

import (
	"bytes"
	"context"
	"encoding/json"
	"entity"
	"fmt"
	"net/http"
	"net/url"
	"services/assert"
	"time"

	"demands/internal/models"

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
	handicap    int64
}

func (*demand) CallRate() int {
	// TODO : implement me
	return 100
}

func (d *demand) Handicap() int64 {
	return d.handicap
}

func (d *demand) Name() string {
	return d.key
}

func (d *demand) hasLimits() bool {
	if d.minuteLimit == 0 &&
		d.hourLimit == 0 &&
		d.dayLimit == 0 &&
		d.weekLimit == 0 &&
		d.monthLimit == 0 {
		return true
	}
	mo, we, da, ho, mi := getCPM(d.key)
	if mo > 0 && mo >= d.monthLimit {
		return false
	}
	if we > 0 && we >= d.weekLimit {
		return false
	}
	if da > 0 && da >= d.dayLimit {
		return false
	}
	if ho > 0 && ho >= d.hourLimit {
		return false
	}
	if mi > 0 && mi >= d.minuteLimit {
		return false
	}
	return true
}

func (d *demand) Provide(ctx context.Context, imp entity.Impression, ch chan map[string]entity.Advertise) {
	defer close(ch)
	if !d.hasLimits() {
		return
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(imp.Raw()); err != nil {
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

	ads := map[string]*restAd{}
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err := dec.Decode(&ads); err != nil {
		logrus.Debug(err)
		return
	}

	adsInter := make(map[string]entity.Advertise, len(ads))
	for _, sl := range imp.Slots() {
		tmp := ads[sl.TrackID()]
		if tmp == nil {
			continue
		}
		// set demand for win resp
		tmp.demand = d
		adsInter[sl.TrackID()] = tmp
	}

	ch <- adsInter
}

func (d *demand) Status(c context.Context, h http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (d *demand) Win(ctx context.Context, id string, cpm int64) {
	incCPM(d.key, cpm)
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
	dm := &demand{
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
		handicap:           d.Handicap,
	}
	dm.createConnection()
	return dm
}
