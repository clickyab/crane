package restful

import (
	"context"
	"encoding/base64"
	"net/http"
	"services/assert"
	"services/eav"
	"strconv"

	core2 "octopus/core"

	"services/safe"

	"octopus/exchange/materialize"
	"services/broker"

	"github.com/fzerorubigd/xmux"
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
			return
		}
		//get from store
		store := eav.NewEavStore(trackID).AllKeys()
		winnerDemand := store["DEMAND"]
		winnerID := store["ID"]
		slotTrack := store["TRACK"]
		AdID := store["ADID"]
		winnerBID := store["BID"]
		IP := store["IP"]
		impTime := store["TIME"]
		supplier := store["SUPPLIER"]
		publisher := store["PUBLISHER"]
		winnerInt, err := strconv.ParseInt(winnerBID, 10, 0)
		if err != nil {
			return
		}
		//set winner
		d, err := core2.GetDemand(winnerDemand)
		if err != nil {
			return
		}
		d.Win(ctx, winnerID, winnerInt)
		broker.Publish(materialize.ShowJob(
			trackID, winnerDemand, slotTrack, AdID, IP, winnerInt, impTime, supplier, publisher,
		))
	})
}

func init() {
	var err error
	data, err = base64.StdEncoding.DecodeString(message)
	assert.Nil(err)
}
