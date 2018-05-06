package native

import (
	"encoding/json"
	"fmt"
	"net/http"

	"clickyab.com/crane/demand/entity"
	"github.com/bsm/openrtb"
	"github.com/bsm/openrtb/native/request"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/framework"
)

func getImps(r *http.Request, count int, pub entity.Publisher, image bool) []openrtb.Impression {
	var (
		res []openrtb.Impression
		sec = secure(r)
	)
	// calculate min cpc and insert in impression ext
	impExt := map[string]interface{}{
		"min_cpc": pub.MinCPC(string(entity.RequestTypeNative)),
	}
	iExt, err := json.Marshal(impExt)
	assert.Nil(err)
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
			},
		}
		if image {
			req.Assets = append(req.Assets, request.Asset{
				ID:       2, // ID 2 is for image
				Required: 1,
				Image: &request.Image{
					TypeID: request.ImageTypeMain,
				},
			})
		}

		bReq, err := json.Marshal(req)
		assert.Nil(err)
		imp := openrtb.Impression{
			ID:       fmt.Sprintf("%d%s%d", pub.ID(), "470", i),
			BidFloor: float64(pub.FloorCPM()),
			Secure:   sec,
			Native: &openrtb.Native{
				Request: bReq,
			},
			Ext: iExt,
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
