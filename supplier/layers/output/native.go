package output

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	"encoding/json"

	"github.com/bsm/openrtb"
	"github.com/bsm/openrtb/native/response"
)

type nativeResp struct {
	Title      string `json:"title"`
	Image      string `json:"image"`
	Impression string `json:"impression"`
	Click      string `json:"click"`
}

const defaultTemplateText string = `
{{range $index, $results := .}}
	{{if isOdd $index}}
		<div>
	{{end}}
	<div>
        <a href='{{.Click}}'>
            <img src='{{.Image}}'/‎>
            <img style="display: none;" src='{{.Impression}}'/‎>
            <div>{{.Title}}</div>
        </a>
    </div>
	{{if isEven $index}}
		</div>
	{{end}}
{{end}}
`

var (
	templateFuncs = template.FuncMap{
		"isEven": isEven,
		"isOdd":  isOdd,
	}
	nativeTemplate = template.Must(template.New("inapp-template").Funcs(templateFuncs).Parse(defaultTemplateText))
)

func isEven(x int) bool {
	return (x+1)%2 == 0
}

func isOdd(x int) bool {
	return (x+1)%2 != 0
}

// RenderNative is the native ad renderer
func RenderNative(_ context.Context, resp *openrtb.BidResponse) ([]byte, error) {
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

	var outputHTML bytes.Buffer
	err := nativeTemplate.Execute(&outputHTML, res)
	if err != nil {
		return nil, err
	}

	r := map[string]string{"html": outputHTML.String()}
	b, err := json.Marshal(r)

	return b, err
}
