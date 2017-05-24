package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"clickyab.com/exchange/services/assert"

	"github.com/Sirupsen/logrus"
)

var client *http.Client

// IP2lData struct
type IP2lData struct {
	CountryShort string `json:"country_short"`
	CountryLong  string `json:"country_long"`
	Region       string `json:"region"`
	City         string `json:"city"`
}

func createConnection() {
	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 30,
		},
		Timeout: 100 * time.Millisecond,
	}
}

// IP2Location return IP2lData from IP
func IP2Location(ip string) IP2lData {
	target := fmt.Sprintf("http://%s/%s", ip2lserver, ip)
	logrus.Debug(target)
	req, err := http.NewRequest("GET", target, nil)
	assert.Nil(err)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Debug(err)
		return IP2lData{
			City:         "-",
			CountryLong:  "-",
			CountryShort: "-",
			Region:       "-",
		}
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	i := IP2lData{}
	err = dec.Decode(&i)
	if err != nil {
		logrus.Debug(err)
		return IP2lData{
			City:         "-",
			CountryLong:  "-",
			CountryShort: "-",
			Region:       "-",
		}
	}

	return i
}

func init() {
	createConnection()
}
