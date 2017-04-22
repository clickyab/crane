package demands

import "octopus/exchange"

type rawPub struct {
	// Name of publisher
	Name string `json:"name"`
	// FloorCPM is the floor cpm for publisher
	FloorCPM int64 `json:"floor_cpm"`
	// SoftFloorCPM is the soft version of floor cpm. if the publisher ahs it, then the system
	// try to use this as floor, but if this is not available, the FloorCPM is used
	SoftFloorCPM int64 `json:"soft_floor_cpm"`
	// Attributes is the generic attribute system
	Attributes map[string]interface{} `json:"attributes"`
}

func getRawPub(in exchange.Publisher) rawPub {
	return rawPub{
		Name:         in.Name(),
		Attributes:   in.Attributes(),
		SoftFloorCPM: in.SoftFloorCPM(),
		FloorCPM:     in.FloorCPM(),
	}
}

type rawLocation struct {
	// Country get the country if available
	Country exchange.Country
	// Province get the province of request if available
	Province exchange.Province
	// LatLon return the latitude longitude if any
	LatLon exchange.LatLon
}

func getRawLocation(in exchange.Location) rawLocation {
	return rawLocation{
		LatLon:   in.LatLon(),
		Province: in.Province(),
		Country:  in.Country(),
	}
}

type rawSlot struct {
	// Size return the primary size of this slot
	Width  int `json:"width"`
	Height int `json:"height"`
	// TrackID is an string for this slot, its a random at first but the value is not changed at all other calls
	TrackID string `json:"track_id"`
}

func getRawSlots(in []exchange.Slot) []rawSlot {
	res := make([]rawSlot, len(in))
	for i := range in {
		res[i] = rawSlot{
			Width:   in[i].Width(),
			Height:  in[i].Height(),
			TrackID: in[i].TrackID(),
		}
	}

	return res
}

type rawImp struct {
	TrackID   string `json:"track_id"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	// Source return the publisher that this client is going into system from that
	Source rawPub `json:"source"`
	// Location of the request
	Location rawLocation `json:"location"`
	// Attributes is the generic attribute system
	Attributes map[string]interface{} `json:"attributes"`
	// Slots is the slot for this request
	Slots []rawSlot `json:"slots"`
	// Category returns category obviously
	Category []exchange.Category `json:"category"`
	// Type return the publisher type
	Type string `json:"type"`
	// Is this publisher accept under floor ads or not ?
	UnderFloor bool `json:"under_floor"`
}

func getRawImpresssion(imp exchange.Impression) interface{} {
	checkIP := func() string {
		if i := imp.IP(); i != nil {
			return i.String()
		}
		return ""
	}
	return rawImp{
		TrackID:    imp.TrackID(),
		IP:         checkIP(),
		UserAgent:  imp.UserAgent(),
		Attributes: imp.Attributes(),
		Category:   imp.Category(),
		Type:       string(imp.Type()),
		UnderFloor: imp.UnderFloor(),
		Slots:      getRawSlots(imp.Slots()),
		Location:   getRawLocation(imp.Location()),
		Source:     getRawPub(imp.Source()),
	}

}
