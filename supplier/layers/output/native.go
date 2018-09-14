package output

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"clickyab.com/crane/openrtb"
	"github.com/bsm/openrtb/native/response"
)

type nativeResp struct {
	Title      string `json:"title"`
	Image      string `json:"image"`
	Impression string `json:"impression"`
	Click      string `json:"click"`
}

// RenderNative is the native ad renderer
func RenderNative(_ context.Context, resp *openrtb.BidResponse, tpl *template.Template) ([]byte, error) {
	var res []nativeResp

	for i := range resp.GetSeatbid() {
		AllBid := resp.GetSeatbid()[i].GetBid()
		if len(AllBid) < 1 {
			continue
		}
		bid := AllBid[0]
		markup := strings.Replace(bid.GetAdm(), "${AUCTION_PRICE}", fmt.Sprint(bid.Price), -1)
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

	var outputHTML bytes.Buffer
	err := tpl.Execute(&outputHTML, res)
	if err != nil {
		return nil, err
	}

	r := map[string]string{"html": outputHTML.String()}
	b, err := json.Marshal(r)

	return b, err
}
