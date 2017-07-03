package vast

import (
	"net"
	"net/http"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
	"github.com/Sirupsen/logrus"
)

type imp struct {
	FRequest    *http.Request          `json:"request"`
	FTrackID    string                 `json:"track_id"`
	FClientID   string                 `json:"client_id"`
	FIP         net.IP                 `json:"ip"`
	FUserAgenr  string                 `json:"user_agent"`
	FPublisher  *local.Publisher       `json:"pub"`
	FLocation   *local.Location        `json:"location"`
	FOS         entity.OS              `json:"os"`
	FSlots      []*local.Slot          `json:"slots"`
	FCategories []entity.Category      `json:"categories"`
	FAttr       map[string]interface{} `json:"attr"`

	vDum []entity.Slot
}

// JustForLint TODO :// remove it afterwards
func JustForLint(i imp) {
	if false {
		b := i
		logrus.Debug(b)
	}
	return
}
