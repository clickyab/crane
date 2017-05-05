package renderer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"octopus/exchange"
)

type dumbAd struct {
	ID      string `json:"id"`
	Winner  int64  `json:"winner"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Code    string `json:"code"`
	Landing string `json:"land"`
}

type restful struct {
	pixelPattern *url.URL
	sup          exchange.Supplier
}

func (rf restful) Render(in map[string]exchange.Advertise, w io.Writer) error {
	res := make(map[string]*dumbAd, len(in))
	for i := range in {
		if in[i] == nil {
			res[i] = nil
			continue
		}

		d := &dumbAd{
			ID:      in[i].ID(),
			Winner:  in[i].WinnerCPM() * int64(100-rf.sup.Share()) / 100,
			Width:   in[i].Width(),
			Height:  in[i].Height(),
			Landing: in[i].Landing(),
		}
		var x url.URL = *rf.pixelPattern
		q := x.Query()
		q.Set("id", i)
		x.RawQuery = q.Encode()

		winURL := in[i].URL()
		win, err := url.Parse(winURL)
		if err == nil {
			q := win.Query()
			q.Set("win", fmt.Sprint(in[i].WinnerCPM()))
			win.RawQuery = q.Encode()
			winURL = win.String()
		}
		// TODO : add the track pixel too x.string()
		d.Code = winURL
		res[i] = d
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
