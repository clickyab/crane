package native

import (
	"net/http"

	"github.com/Sirupsen/logrus"

	"net"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
)

type imp struct {
	FRequest    *http.Request          `json:"request"`
	FTrackID    string                 `json:"track_id"`
	FClientID   string                 `json:"client_id"`
	FIP         net.IP                 `json:"ip"`
	FUA         string                 `json:"user_agent"`
	FPub        *local.Publisher       `json:"pub"`
	FLocation   entity.Location        `json:"location"`
	FOS         entity.OS              `json:"os"`
	FSlots      []*local.Slot          `json:"slots"`
	FCategories []entity.Category      `json:"categories"`
	FAttr       map[string]interface{} `json:"attr"`

	nDum   []entity.Slot
	latlon entity.LatLon
}

// JustForLint TODO :// remove it afterwards
func JustForLint(i imp) {
	if false {
		b := i.latlon
		logrus.Debug(b)
		i.extractData()
	}
	return
}
