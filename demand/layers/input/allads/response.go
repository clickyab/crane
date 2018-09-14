package allads

import "clickyab.com/crane/demand/entity"

type response struct {
	ExceedAds     []responseAds `json:"exceed_ads"`
	UnderfloorAds []responseAds `json:"underfloor_ads"`
	FilteredAds   []responseAds `json:"filtered_ads"`
}

type responseAds struct {
	ID            int32          `json:"id"`
	CampaignID    int32          `json:"campaign_id"`
	MaxBid        int32          `json:"max_bid"`
	Errors        []string       `json:"errors,omitempty"`
	CampaignName  string         `json:"campaign_name,omitempty"`
	Type          entity.AdType  `json:"type,omitempty"`
	TargetURL     string         `json:"target_url,omitempty"`
	Size          int32          `json:"size,omitempty"`
	Frequency     int32          `json:"frequency,omitempty"`
	CalculatedCPC float64        `json:"calculated_cpc,omitempty"`
	CalculatedCPM float64        `json:"calculated_cpm,omitempty"`
	CalculatedCTR float64        `json:"calculated_ctr,omitempty"`
	Assets        []entity.Asset `json:"assets"`
}
