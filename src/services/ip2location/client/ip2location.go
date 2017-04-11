package client

import (
	"fmt"
	"net/http"
	"time"

	"assert"
	"encoding/json"
)

var client *http.Client

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

func IP2Location(ip string) IP2lData {
	target := fmt.Sprintf("http://%s/%s", ip2lserver, ip)
	req, err := http.NewRequest("GET", target, nil)
	assert.Nil(err)
	resp, err := client.Do(req)
	if err != nil {
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
