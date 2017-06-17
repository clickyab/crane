package restful

import (
	"context"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/eav"

	core2 "clickyab.com/exchange/octopus/core"

	"github.com/clickyab/services/safe"

	"clickyab.com/exchange/octopus/exchange/materialize"
	"github.com/clickyab/services/broker"

	"github.com/Sirupsen/logrus"
	"github.com/rs/xmux"
)

const message = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEUAAACnej3aAAAAAXRSTlMAQObYZgAAAApJREFUCNdjYAAAAAIAAeIhvDMAAAAASUVORK5CYII="

var data []byte

// TrackPixel is a route to handle track pixel event
func TrackPixel(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Write(data)
	safe.GoRoutine(func() {
		// TODO demand name must be solid
		demand := xmux.Param(ctx, "demand")
		trackID := xmux.Param(ctx, "trackID")
		if trackID == "" || demand == "" {
			logrus.Debug("both track id and demand are empty")
			return
		}
		//get from store
		store := eav.NewEavStore("PIXEL_" + trackID).AllKeys()
		winnerDemand := store["DEMAND"]
		if winnerDemand != demand {
			logrus.Debugf("stored demand `%s`!=request demand `%s`", winnerDemand, demand)
			return
		}
		winnerID := store["ID"]
		slotTrack := store["TRACK"]
		AdID := store["ADID"]
		winnerBID := store["BID"]
		IP := store["IP"]
		impTime := store["TIME"]
		supplier := store["SUPPLIER"]
		publisher := store["PUBLISHER"]
		profit := store["PROFIT"]
		profitInt, err := strconv.ParseInt(profit, 10, 0)
		if err != nil {
			logrus.Debugf("profit key is not integer : %s", err)
			return
		}
		winnerInt, err := strconv.ParseInt(winnerBID, 10, 0)
		if err != nil {
			logrus.Debugf("winner key is not integer : %s", err)
			return
		}
		//set winner
		d, err := core2.GetDemand(winnerDemand)
		if err != nil {
			logrus.Debugf("can not find winner demand `%s` : %s", winnerDemand, err)
			return
		}
		broker.Publish(materialize.ShowJob(
			trackID, winnerDemand, slotTrack, AdID, IP, winnerInt, impTime, supplier, publisher, profitInt,
		))
		d.Win(ctx, winnerID, winnerInt)
	})
}

func init() {
	var err error
	data, err = base64.StdEncoding.DecodeString(message)
	assert.Nil(err)
}
