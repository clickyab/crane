package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/influx"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/ip2location"
	"github.com/clickyab/services/safe"
	"github.com/clickyab/services/shell"

	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"github.com/mssola/user_agent"
	"gopkg.in/fzerorubigd/onion.v3"
)

// Point is the impression point
type Point struct {
	CampaignID    string `json:"cp_id"`
	PublisherID   string `json:"pub_id"`
	PublisherType string `json:"publisher_type"`
	WinnerBID     int64  `json:"winner_bid"`
	UserAgent     string `json:"user_agent"`
	IP            string `json:"ip"`
}

var (
	port onion.String
)

func main() {
	config.Initialize("clickyab", "collector", "CC")
	defer initializer.Initialize()()

	mux := xmux.New()
	mux.POST("/", xhandler.HandlerFuncC(collect))

	srv := &http.Server{Addr: port.String(), Handler: xhandler.New(context.Background(), mux)}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Debug(err)
		}
	}()

	shell.WaitExitSignal()
	s, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	srv.Shutdown(s)

}

func collect(_ context.Context, w http.ResponseWriter, r *http.Request) {
	safe.GoRoutine(func() {
		data, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		p := Point{}
		fmt.Println(string(data))
		err := json.Unmarshal(data, &p)
		if err != nil {
			logrus.Error(err)
			return
		}

		rec := ip2location.GetAll(p.IP)
		ua := user_agent.New(p.UserAgent)
		influx.AddPoint(
			"impression",
			map[string]string{
				"campaign":       fmt.Sprint(p.CampaignID),
				"publisher":      p.PublisherID,
				"publisher_type": p.PublisherType,
				"country":        rec.CountryLong,
				"province":       rec.Region,
				"city":           rec.City,
				"mobile":         fmt.Sprint(ua.Mobile()),
				"platform":       fmt.Sprint(ua.Platform()),
			},
			map[string]interface{}{
				"count":      1,
				"winner_bid": p.WinnerBID,
			},
			time.Now(),
		)
	})
	w.WriteHeader(204)
}

func init() {
	p := os.Getenv("PORT")
	if p == "" {
		p = "8878"
	}

	port = config.RegisterString("collector.listen", fmt.Sprintf(":%s", p), "the port for collector api endpoint")
}
