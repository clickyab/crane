package url

import (
	"fmt"

	"clickyab.com/crane/crane/entity"

	"encoding/base64"
	"net/http"
	"time"

	"strconv"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/eav"
	"github.com/clickyab/services/safe"
)

var expire = config.GetDurationDefault("crane.output.url_expire", time.Hour*72)
var fastclick = config.GetDurationDefault("crane.output.fastclick", time.Second*5)

var host = config.GetStringDefault("crane.output.click_host", "clickyab.com")

const prefix = "click"

const (
	impressionTrackID = "a"
	publisher         = "b"
	clientID          = "c"
	supplier          = "d"
	ip                = "e"
	userAgent         = "f"
	slotTrackID       = "g"
	winnerBID         = "h"
	advertiseID       = "i"
	genTime           = "j"
)

type click struct {
}

// ClickURL generate url for click event
func (click) ClickURL(s entity.Slot, m entity.Context) string {
	// pattern= http://server.com/publisher/slotTrackID/impTrackID/adID/adType/?t=target(base64)
	targetURL := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(s.WinnerAdvertise().TargetURL()))
	url := fmt.Sprintf("http://%s/%s/%s/%s/%s/%s?t=%s", host, m.Publisher().Name(), s.TrackID(), m.TrackID(), s.WinnerAdvertise().ID(), s.WinnerAdvertise().Type(), targetURL)
	safe.GoRoutine(func() {
		eav.NewEavStore(keyGen(m.TrackID(), s.TrackID())).
			SetSubKey(impressionTrackID, m.TrackID()).
			SetSubKey(publisher, m.Publisher().Name()).
			SetSubKey(clientID, m.ClientID()).
			SetSubKey(supplier, m.Publisher().Supplier()).
			SetSubKey(ip, m.IP().String()).
			SetSubKey(userAgent, m.UserAgent()).
			SetSubKey(slotTrackID, s.TrackID()).
			SetSubKey(winnerBID, strconv.FormatInt(s.WinnerAdvertise().WinnerBID(), 36)).
			SetSubKey(advertiseID, s.WinnerAdvertise().ID()).
			SetSubKey(genTime, strconv.FormatInt(time.Now().Unix(), 36)).
			Save(expire)
	})
	return url

}

// ClickHandler process click url
func ClickHandler(c Data, w http.ResponseWriter) {

	// url e.g. /publisher/slotTrackID/impTrackID/adID/adType/?t=target(base64)

	u, e := base64.StdEncoding.DecodeString(c.Target)
	if e != nil {
		w.Write([]byte(backTemplate))
		w.WriteHeader(http.StatusNotFound)
	} else {
		// add referral metadata to this map
		m := map[string]interface{}{
			"publisher":  c.Publisher,
			"impression": c.ImpressionTrackID,
			"cid":        c.ClientID,
		}
		w.Header().Set("location", addMeta(string(u), m))
		w.WriteHeader(http.StatusPermanentRedirect)
	}

	assert.NotNil(worker)

	safe.GoRoutine(func() {

		ki := eav.NewEavStore(keyGen(c.ImpressionTrackID, c.SlotTrackID))
		v := ki.AllKeys()

		if len(v) == 0 {
			worker <- Data{

				ImpressionTrackID: c.ImpressionTrackID,
				ClientID:          c.ClientID,
				IP:                c.IP,
				UserAgent:         c.UserAgent,
				SlotTrackID:       c.SlotTrackID,
				AdvertiseID:       c.AdvertiseID,
				ClickTime:         time.Now(),
				Status:            entity.SuspNoAdFound,
			}
			return
		}

		w, _ := strconv.ParseInt(v[winnerBID], 36, 64)
		a, _ := strconv.ParseInt(v[advertiseID], 36, 64)
		vt, _ := strconv.ParseInt(v[genTime], 36, 64)
		t := time.Unix(vt, 0)
		d := Data{
			ImpressionTrackID: v[impressionTrackID],
			Publisher:         v[publisher],
			ClientID:          v[clientID],
			Supplier:          v[supplier],
			IP:                v[ip],
			UserAgent:         v[userAgent],
			SlotTrackID:       v[slotTrackID],
			WinnerBID:         w,
			AdvertiseID:       a,
			GenTime:           t,
			ClickTime:         time.Now(),
		}

		if d.UserAgent != c.UserAgent {
			d.Status = entity.SuspUAMismatch
			worker <- d
			return

		}
		if d.IP != c.IP {
			d.Status = entity.SuspIPMismatch
			worker <- d
			return
		}
		if fastclick > time.Duration(d.ClickTime.Unix()-d.GenTime.Unix()) {
			d.Status = entity.SuspFastClick
			worker <- d
			return
		}
		d.Status = entity.SuspSuccessful
		worker <- d
	})

}

const backTemplate = `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title></title>
    <script>
    setTimeout(function () {
        window.history.backTemplate()
    },1500)
    </script>
    <style>
body {
  background: #dddddd;
  color: #222;
  direction: rtl;
}
body .txt {
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  margin: 0 auto;
  padding: 25px;
  height: 100vh;
  max-width: 700px;
}
    </style>
</head>
<body >

  <h1 class="txt">
   صفحه مورد نظر یافت نشد!
  </h1>
</body>
</html>`
