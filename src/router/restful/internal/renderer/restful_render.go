package renderer

import (
	"context"
	"encoding/json"
	"entity"
	"fmt"
	"io"
	"net/url"
)

type dumbAd struct {
	ID     string `json:"id"`
	MaxCPM int64  `json:"max_cpm"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Code   string `json:"code"`
}

type restful struct {
	pixelPattern *url.URL
}

func (rf restful) Render(ctx context.Context, in map[string]entity.Advertise, w io.Writer) error {
	res := make(map[string]*dumbAd, len(in))
	for i := range in {
		if in[i] == nil {
			res[i] = nil
			continue
		}

		d := &dumbAd{
			ID:     in[i].ID(),
			MaxCPM: in[i].MaxCPM(),
			Width:  in[i].Width(),
			Height: in[i].Height(),
		}
		var x url.URL = *rf.pixelPattern
		q := x.Query()
		q.Set("id", i)
		x.RawQuery = q.Encode()
		d.Code = fmt.Sprintf("TODO , tracker code is %s, the actual route is %s", x.String(), in[i].URL())
		res[i] = d
	}

	enc := json.NewEncoder(w)
	return enc.Encode(res)
}

// NewRestfulRenderer return a restful renderer
func NewRestfulRenderer(pixel *url.URL) entity.Renderer {
	return &restful{
		pixelPattern: pixel,
	}
}
