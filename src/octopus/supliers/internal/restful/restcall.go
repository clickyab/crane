package restful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"octopus/exchange"
	"strings"
)

type requestBody struct {
	IP         string              `json:"ip"`
	Publisher  *restPublisher      `json:"publisher"`
	Categories []exchange.Category `json:"categories"`
	Type       string              `json:"type"`
	UnderFloor bool
	App        struct {
		OSVersion  string `json:"os_version,omitempty"`
		Operator   string `json:"operator,omitempty"`
		Brand      string `json:"brand,omitempty"`
		Model      string `json:"model,omitempty"`
		Language   string `json:"language,omitempty"`
		Network    string `json:"network,omitempty"`
		OSIdentity string `json:"os_identity,omitempty"`
		MCC        int64  `json:"mcc,omitempty"`
		MNC        int64  `json:"mnc,omitempty"`
		LAC        int64  `json:"lac,omitempty"`
		CID        int64  `json:"cid,omitempty"`
		UserAgent  string `json:"user_agent,omitempty"`
	} `json:"app,omitempty"`
	Web struct {
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

// GetImpression try to create an impression object from a request
func GetImpression(sup exchange.Supplier, r *http.Request) (exchange.Impression, error) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	rb := requestBody{}
	err := dec.Decode(&rb)
	if err != nil {
		return nil, err
	}

	var res *impressionRest
	switch strings.ToLower(rb.Type) {
	case "app":
		res, err = newImpressionFromAppRequest(sup, &rb)
	case "web":
		res, err = newImpressionFromWebRequest(sup, &rb)
	case "vast":
		res, err = newImpressionFromVastRequest(sup, &rb)
	default:
		err = fmt.Errorf("type is not supported: %s", rb.Type)
	}

	if err != nil {
		return nil, err
	}

	// Hidden profit is here. the floor and soft floor are rising here
	share := int64(100 + res.Pub.sup.Share())
	res.Pub.PubFloorCPM = (res.Pub.PubFloorCPM * share) / 100
	if res.Pub.PubFloorCPM == 0 {
		res.Pub.PubFloorCPM = (sup.FloorCPM() * share) / 100
	}
	res.Pub.PubSoftFloorCPM = (res.Pub.PubSoftFloorCPM * share) / 100
	if res.Pub.PubSoftFloorCPM == 0 {
		res.Pub.PubSoftFloorCPM = (sup.SoftFloorCPM() * share) / 100
	}
	return res, nil
}
