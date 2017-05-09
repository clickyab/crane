package renderer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"octopus/exchange"
	"services/config"
)

type dumbAd struct {
	TrackID   string `json:"track_id"`
	AdTrackID string `json:"ad_track_id"`
	Winner    int64  `json:"winner"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Code      string `json:"code"`
	IsFilled  bool   `json:"is_filled"`
	Landing   string `json:"landing"`
}

type restful struct {
	pixelPattern *url.URL
	sup          exchange.Supplier
}

func (rf restful) Render(imp exchange.Impression, in map[string]exchange.Advertise, w io.Writer) error {
	res := make([]*dumbAd, 0)
	slots := imp.Slots()
	for k := range slots {
		slotTrackID := slots[k].TrackID()
		if in[slotTrackID] == nil {
			ctx := templateContext{
				URL: slots[k].Fallback(),
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

		host := config.GetStringDefault("exchange.host.name", "localhost:3412")
		trackURL := fmt.Sprintf(`%s/pixel/%s/%s`, host, in[slotTrackID].Demand().Name(), in[slotTrackID].TrackID())
		ctx := templateContext{
			URL:      winURL,
			IsFilled: true,
			Landing:  in[slotTrackID].Landing(),
			Pixel:    trackURL,
		}
		d.Code = renderTemplate(ctx)
		res = append(res, d)
	}

	enc := json.NewEncoder(w)
	return enc.Encode(res)
}

// NewRestfulRenderer return a restful renderer
func NewRestfulRenderer(sup exchange.Supplier, pixel *url.URL) exchange.Renderer {
	return &restful{
		pixelPattern: pixel,
		sup:          sup,
	}
}
