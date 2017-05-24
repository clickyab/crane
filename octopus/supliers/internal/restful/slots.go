package restful

type slotRest struct {
	W           int    `json:"width"`
	H           int    `json:"height"`
	TID         string `json:"track_id"`
	FallbackURL string `json:"fallback_url"`
}

func (sr slotRest) Fallback() string {
	return sr.FallbackURL
}

func (sr slotRest) Width() int {
	return sr.W
}

func (sr slotRest) Height() int {
	return sr.H
}

func (sr slotRest) TrackID() string {
	return sr.TID
}
