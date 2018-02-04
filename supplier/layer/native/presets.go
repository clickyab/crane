package native

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bsm/openrtb"
	"github.com/bsm/openrtb/native/request"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/framework"
)

func getImps(r *http.Request, count int, pubID, pubFloor int64) []openrtb.Impression {
	var (
		res []openrtb.Impression
		sec = secure(r)
	)
	for i := 1; i <= count; i++ {
		// make request
		req := request.Request{
			Assets: []request.Asset{
				{
					ID:       1, // ID 1 is for text
					Required: 1,
					Title: &request.Title{
						Length: nativeMaxTitleLen.Int(),
					},
				},
				{
					ID:       2, // ID 2 is for image
					Required: 1,
					Image: &request.Image{
						TypeID: request.ImageTypeMain,
					},
				},
			},
		}
		bReq, err := json.Marshal(req)
		assert.Nil(err)
		imp := openrtb.Impression{
			ID:       fmt.Sprintf("%d%s%d", pubID, "470", i),
			BidFloor: float64(pubFloor),
			Secure:   sec,
			Native: &openrtb.Native{
				Request: bReq,
			},
		}
		res = append(res, imp)
	}
	return res
}

// secure check openrtb protocol (http/https)
func secure(r *http.Request) openrtb.NumberOrString {
	if framework.Scheme(r) == "https" {
		return openrtb.NumberOrString(1)
	}
	return openrtb.NumberOrString(0)
}
