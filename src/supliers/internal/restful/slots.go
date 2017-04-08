package restful

type slotRest struct {
	W          int    `json:"width"`
	H          int    `json:"height"`
	MaximumCPM int64  `json:"max_cpm"`
	TID        string `json:"track_id"`
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

func (sr slotRest) MaxCPM() int64 {
	return sr.MaximumCPM
}
