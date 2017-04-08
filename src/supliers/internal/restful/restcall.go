package restful

import (
	"encoding/json"
	"entity"
	"fmt"
	"net/http"
	"strings"
)

type requestBody struct {
	IP         string            `json:"ip"`
	Publisher  restPublisher     `json:"publisher"`
	Categories []entity.Category `json:"categories"`
	Type       string            `json:"type"`
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

	Slots []slotRest `json:"slots"`
}

// GetImpression try to create an impression object from a request
func GetImpression(sup entity.Supplier, r *http.Request) (entity.Impression, error) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	rb := requestBody{}
	err := dec.Decode(&rb)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(rb.Type) {
	case "app":
		return newImpressionFromAppRequest(sup, &rb)
	case "web":
		return newImpressionFromWebRequest(sup, &rb)
	case "vast":
		return newImpressionFromVastRequest(sup, &rb)
	default:
		return nil, fmt.Errorf("type is not supported: %s", rb.Type)
	}
}
