package output

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"

	"clickyab.com/crane/openrtb/v2.5"
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
		d := nativeResp{
			Impression: bid.GetAdmNative().Imptrackers[0],
			Click:      bid.GetAdmNative().Link.Url,
		}
		for i := range bid.GetAdmNative().GetAssets() {
			if bid.GetAdmNative().GetAssets()[i].Id == 1 {
				d.Title = bid.GetAdmNative().GetAssets()[i].GetTitle().Text
			}
			if bid.GetAdmNative().GetAssets()[i].Id == 2 {
				d.Image = bid.GetAdmNative().GetAssets()[i].GetImg().Url
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
