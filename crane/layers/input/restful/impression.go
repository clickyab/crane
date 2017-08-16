package restful

import (
	"net"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
)

type rawImpression struct {
	FTrackID   string `json:"track_id"`
	FIP        string `json:"ip"`
	FProtocol  string `json:"scheme"`
	FUserAgent string `json:"user_agent"`

	FSource struct {
		Name         string                 `json:"name"`
		Supplier     string                 `json:"supplier"`
		FloorCPM     int64                  `json:"floor_cpm"`
		SoftFloorCPM int64                  `json:"soft_floor_cpm"`
		Attributes   map[string]interface{} `json:"attributes"`
	} `json:"source"`
	FLocation struct {
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
			Long  float64 `json:"long"`
		} `json:"latlon"`
	} `json:"location"`
	FAttributes map[string]string `json:"attributes"`
	FSlots      []local.Slot      `json:"slots"`

	FCategory []entity.Category `json:"category"`

	FPlatform    string `json:"platform"`
	FUnderfloor  bool   `json:"underfloor"`
	FSessionKey  string `json:"page_track_id"`
	FUserTrackID string `json:"user_track_id"`

	publisher entity.Publisher
}

func (r *rawImpression) ID() string {
	panic("implement me")
}

func (r *rawImpression) Width() int {
	panic("implement me")
}

func (r *rawImpression) Height() int {
	panic("implement me")
}

func (r *rawImpression) SetSlotCTR(float64) {
	panic("implement me")
}

func (r *rawImpression) SlotCTR() float64 {
	panic("implement me")
}

func (r *rawImpression) SetWinnerAdvertise(entity.Advertise) {
	panic("implement me")
}

func (r *rawImpression) WinnerAdvertise() entity.Advertise {
	panic("implement me")
}

func (r *rawImpression) SetShowURL(string) {
	panic("implement me")
}

func (r *rawImpression) ShowURL() string {
	panic("implement me")
}

func (r *rawImpression) IsSizeAllowed(int, int) bool {
	panic("implement me")
}

func (r *rawImpression) Attribute() map[string]interface{} {
	panic("implement me")
}

func (r *rawImpression) IP() net.IP {
	return net.ParseIP(r.FIP)
}

func (r *rawImpression) OS() entity.OS {
	return local.OSFromAgent(r.FUserAgent)
}

func (r *rawImpression) ClientID() string {
	return r.FUserTrackID
}

func (r *rawImpression) Protocol() string {
	return r.FProtocol
}

func (r *rawImpression) UserAgent() string {
	return r.FUserAgent
}

func (r *rawImpression) Location() entity.Location {
	l := r.FLocation
	res := &local.Location{
		FCountry: entity.Country{
			Name:  l.Country.Name,
			Valid: l.Country.Valid,
			ISO:   l.Country.ISO,
		},
		FProvince: entity.Province{
			Name:  l.Province.Name,
			Valid: l.Province.Valid,
		},
		FLatLon: entity.LatLon{
			Lat:   l.LatLon.Lat,
			Valid: l.LatLon.Valid,
			Lon:   l.LatLon.Long,
		},
	}
	return res
}

func (r *rawImpression) Attributes() map[string]string {
	return r.FAttributes
}

func (r *rawImpression) TrackID() string {
	return r.FTrackID
}

func (r *rawImpression) Publisher() entity.Publisher {
	return r.publisher
}

func (r *rawImpression) Slots() []entity.Slot {
	a := make([]entity.Slot, 0)
	for i := range r.FSlots {
		a = append(a, &r.FSlots[i])
	}
	return a
}

func (r *rawImpression) Category() []entity.Category {
	return r.FCategory
}
