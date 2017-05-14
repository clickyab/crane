package renderer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"octopus/exchange"
	"services/config"
)

var host = config.RegisterString("octopus.host.name", "exchange-dev.3rdad.com", "the exchange root")

type dumbAd struct {
	TrackID   string `json:"track_id" structs:"track_id"`
	AdTrackID string `json:"ad_track_id" structs:"ad_track_id"`
	Winner    int64  `json:"winner" structs:"winner"`
	Width     int    `json:"width" structs:"width"`
	Height    int    `json:"height" structs:"height"`
	Code      string `json:"code" structs:"code"`
	IsFilled  bool   `json:"is_filled" structs:"is_filled"`
	Landing   string `json:"landing" structs:"landing"`
}

type restful struct {
	pixelPattern *url.URL
	sup          exchange.Supplier
}

func (rf restful) Render(imp exchange.Impression, in map[string]exchange.Advertise, w http.ResponseWriter) error {
	res := make([]*dumbAd, 0)
	slots := imp.Slots()
	for k := range slots {
		slotTrackID := slots[k].TrackID()
		if in[slotTrackID] == nil {
			ctx := templateContext{
				URL:    slots[k].Fallback(),
				Width:  slots[k].Width(),
				Height: slots[k].Height(),
			}

			res = append(res, &dumbAd{
				Code:    renderTemplate(ctx),
				TrackID: slotTrackID,
			})
			continue
		}

		d := &dumbAd{
			TrackID:   slotTrackID,
			AdTrackID: in[slotTrackID].TrackID(),
			Winner:    in[slotTrackID].WinnerCPM() * int64(100-rf.sup.Share()) / 100,
			Width:     in[slotTrackID].Width(),
			Height:    in[slotTrackID].Height(),
			Landing:   in[slotTrackID].Landing(),
			IsFilled:  true,
		}
		x := *rf.pixelPattern
		q := x.Query()
		q.Set("id", slotTrackID)
		x.RawQuery = q.Encode()

		winURL := in[slotTrackID].URL()
		win, err := url.Parse(winURL)
		if err == nil {
			q := win.Query()
			q.Set("win", fmt.Sprint(in[slotTrackID].WinnerCPM()))
			win.RawQuery = q.Encode()
			winURL = win.String()
		}

		// TODO : fix the scheme
		trackURL := &url.URL{
			Scheme: "http",
			Host:   *host,
			Path:   fmt.Sprintf(`/pixel/%s/%s`, in[slotTrackID].Demand().Name(), in[slotTrackID].TrackID()),
		}

		ctx := templateContext{
			URL:      winURL,
			IsFilled: true,
			Landing:  in[slotTrackID].Landing(),
			Pixel:    trackURL.String(),
			Width:    slots[k].Width(),
			Height:   slots[k].Height(),
		}
		d.Code = renderTemplate(ctx)
		res = append(res, d)
	}

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return enc.Encode(res)
}

// NewRestfulRenderer return a restful renderer
func NewRestfulRenderer(sup exchange.Supplier, pixel *url.URL) exchange.Renderer {
	return &restful{
		pixelPattern: pixel,
		sup:          sup,
	}
}
