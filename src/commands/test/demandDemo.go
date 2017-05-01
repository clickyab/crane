package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func demandDemo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var asd []byte
	r.Body.Read(asd)
	temp := &payload{}
	err := decoder.Decode(temp)
	if false {
		println(temp.Source.Name)
		println(temp.Source.SoftFloorCPM)
		println(temp.Location.Country.Name)
		println(temp.Location.Country.ISO)
		println(temp.Location.Country.Valid)
		println(temp.Location.Province.Name)
		println(temp.Location.Province.Valid)
		println(temp.Location.LatLon.Valid)
		println(temp.Location.LatLon.Lat)
		println(temp.Location.LatLon.long)
		println(temp.Source.Attributes)
		println(temp.Slots[0].TrackID)
	}
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	var res response
	for i := range temp.Slots {
		a := singleResponse{
			ID:      "1",
			MaxCPM:  temp.Source.FloorCPM + 1,
			Width:   temp.Slots[i].Width,
			Height:  temp.Slots[i].Height,
			URL:     fmt.Sprintf("http://a.clickyab.com/ads/?a=4471405272967&width=%d&height=%d&slot=71634138754&domainname=p30download.com&eventpage=416310534&loc=http%3A%2F%2Fp30download.com%2Fagahi%2Fplan%2Fa1i.php&ref=http%3A%2F%2Fp30download.com%2F&adcount=1", temp.Slots[i].Width, temp.Slots[i].Height),
			Landing: "clickyab.com",
		}
		res = append(res, a)
	}

	b, _ := json.Marshal(res)
	fmt.Println(string(b))

	encoder := json.NewEncoder(w)
	err = encoder.Encode(res)
	if err != nil {
		panic("asd")
	}

}

type payload struct {
	TrackID   string `json:"track_id"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`

	Source struct {
		Name         string                 `json:"name"`
		FloorCPM     int                    `json:"floor_cpm"`
		SoftFloorCPM int                    `json:"soft_floor_cpm"`
		Attributes   map[string]interface{} `json:"attributes"`
	} `json:"source"`

	Location struct {
		Country struct {
			Valid bool   `json:"valid"`
			Name  string `json:"name"`
			ISO   string `json:"iso"`
		} `json:"country"`
		Province struct {
			Valid bool   `json:"valid"`
			Name  string `json:"name"`
		} `json:"province"`
		LatLon struct {
			Valid bool    `json:"valid"`
			Lat   float64 `json:"lat"`
			long  float64 `json:"long"`
		} `json:"latlon"`
	} `json:"location"`

	Attributes map[string]interface{} `json:"attributes"`
	Slots      []struct {
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		TrackID string `json:"track_id"`
	} `json:"slots"`

	Category []string `json:"category"`

	Platform   string `json:"platform"`
	Underfloor bool   `json:"underfloor"`
}

type singleResponse struct {
	ID      string `json:"id"`
	MaxCPM  int    `json:"max_cpm"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	URL     string `json:"url"`
	Landing string `json:"landing"`
}

type response []singleResponse
