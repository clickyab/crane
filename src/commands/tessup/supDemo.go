package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

type categories string

func getSupplierDemo(w http.ResponseWriter, r *http.Request) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(dir)
	t, err := template.ParseFiles("../src/commands/tessup/index.html")
	if err != nil {
		return
	}
	logrus.Debug(t.Execute(w, nil))
}

type dumbAd struct {
	ID      string `json:"id"`
	Winner  int64  `json:"winner"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Code    string `json:"code"`
	Landing string `json:"land"`
}

func postSupplierDemo(w http.ResponseWriter, r *http.Request) {
	var respon map[string]dumbAd
	var cats []categories
	var resSlots = make([]*slotRest, 0)
	if r.Method != "POST" {
		return
	}
	r.ParseForm()
	data := r.Form
	rType := data["type"][0]
	refferer := data["refferer"][0]
	parent := data["parent"][0]
	userAgent := data["user_agent"][0]
	res := requestBody{}
	for i := range data["width[]"] {
		if data["width[]"][i] != "" {
			wi, _ := strconv.ParseInt(data["width[]"][i], 10, 0)
			he, _ := strconv.ParseInt(data["height[]"][i], 10, 0)
			resSlots = append(resSlots, &slotRest{
				W:   int(wi),
				H:   int(he),
				TID: data["track[]"][i],
			})
		}

	}
	if data["categories"][0] != "" {
		category := strings.Split(data["categories"][0], ",")
		for i := range category {
			cats = append(cats, categories(category[i]))
		}

	}
	floorCPM, _ := strconv.ParseInt(data["floor_cpm"][0], 10, 0)
	softFloor, _ := strconv.ParseInt(data["soft_floor"][0], 10, 0)
	res.UnderFloor = true
	res.IP = data["ip"][0]
	res.Categories = cats
	res.Publisher = &restPublisher{
		PubFloorCPM:     floorCPM,
		PubName:         data["publisher_name"][0],
		PubSoftFloorCPM: softFloor,
	}
	res.IP = "46.209.239.51"
	res.Type = rType
	switch rType {
	case "web":
		res.Web.Referrer = refferer
		res.Web.Parent = parent
		res.Web.UserAgent = userAgent
	case "vast":
		res.Vast.Referrer = refferer
		res.Vast.Parent = parent
		res.Vast.UserAgent = userAgent
	}
	res.Slots = resSlots
	resData, err := json.Marshal(res)
	if err != nil {

	}
	logrus.Warn(string(resData))
	request, err := http.NewRequest("POST", "http://127.0.0.1:8090/get/randomhash", bytes.NewBuffer(resData))
	if err != nil {
		return
	}
	request.Header.Add("X-REAL-IP", "46.209.239.50")
	fmt.Println(err)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(responseData, &respon)
	logrus.Warn(err)
	logrus.Warnf("%+v", respon)
	t, err := template.ParseFiles("../src/commands/tessup/show.tmpl")
	w.Header().Set("Content-Type", "text/html")
	logrus.Warn(t.Execute(w, respon))
	fmt.Println(err)
}

type supplier struct {
	Name            string   `json:"name"`
	FloorCPM        int64    `json:"floor_cpm"`
	SoftFloorCPM    int64    `json:"soft_floor_cpm"`
	ExcludedDemands []string `json:"excluded_demands"`
	Share           int      `json:"share"`
}

type slotRest struct {
	W   int    `json:"width"`
	H   int    `json:"height"`
	TID string `json:"track_id"`
}

type restPublisher struct {
	PubName         string `json:"name"`
	PubFloorCPM     int64  `json:"floor_cpm"`
	PubSoftFloorCPM int64  `json:"soft_floor_cpm"`

	Sup supplier `json:"sup"`
}

type requestBody struct {
	IP         string         `json:"ip"`
	Publisher  *restPublisher `json:"publisher"`
	Categories []categories   `json:"categories"`
	Type       string         `json:"type"`
	UnderFloor bool
	Web        struct {
		Referrer  string `json:"referrer,omitempty"`
		Parent    string `json:"parent,omitempty"`
		UserAgent string `json:"user_agent,omitempty"`
	} `json:"web,omitempty"`
	Vast struct {
		Referrer  string `json:"referrer,omitempty"`
		Parent    string `json:"parent,omitempty"`
		UserAgent string `json:"user_agent,omitempty"`
	} `json:"vast,omitempty"`

	Slots []*slotRest `json:"slots"`
}
