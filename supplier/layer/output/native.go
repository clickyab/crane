package output

import (
	"context"
	"fmt"
	"io"
	"strings"

	"encoding/json"

	"github.com/bsm/openrtb"
	"github.com/bsm/openrtb/native/response"
	"github.com/clickyab/services/assert"
)

type nativeResp struct {
	Title      string `json:"title"`
	Image      string `json:"image"`
	Impression string `json:"impression"`
	Click      string `json:"click"`
}

// RenderNative is the native ad renderer
func RenderNative(_ context.Context, w io.Writer, resp *openrtb.BidResponse) error {
	var res []nativeResp

	for i := range resp.SeatBid {
		AllBid := resp.SeatBid[i].Bid
		if len(AllBid) < 1 {
			continue
		}
		bid := AllBid[0]
		markup := strings.Replace(bid.AdMarkup, "${AUCTION_PRICE}", fmt.Sprint(bid.Price), -1)
		bs := response.Response{}
		err := json.Unmarshal([]byte(markup), &bs)
		if err != nil {
			continue
		}
		d := nativeResp{
			Impression: bs.ImpTrackers[0],
			Click:      bs.Link.URL,
		}
		for i := range bs.Assets {
			if bs.Assets[i].ID == 1 {
				d.Title = bs.Assets[i].Title.Text
			}
			if bs.Assets[i].ID == 2 {
				d.Image = bs.Assets[i].Image.URL
			}
		}
		res = append(res, d)
	}

	b, err := json.Marshal(res)
	assert.Nil(err)
	_, _ = w.Write(b)
	return nil
}
