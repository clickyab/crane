package native

import (
	"fmt"
	"net/http"

	"clickyab.com/crane/demand/entity"
	"clickyab.com/crane/openrtb"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
)

func getImps(r *http.Request, count int, pub entity.Publisher, image bool) []*openrtb.Imp {
	var (
		res []*openrtb.Imp
		sec = framework.Scheme(r) == "https"
	)
	for i := 1; i <= count; i++ {
		// make request
		req := &openrtb.Native_RequestNative{
			RequestNative: &openrtb.NativeRequest{
				Assets: []*openrtb.NativeRequest_Asset{
					{
						Id:       1,
						Required: true,
						AssetOneof: &openrtb.NativeRequest_Asset_Title_{
							Title: &openrtb.NativeRequest_Asset_Title{
								Len: int32(nativeMaxTitleLen.Int()),
							},
						},
					},
				},
			},
		}
		if image {
			req.RequestNative.Assets = append(req.RequestNative.Assets, &openrtb.NativeRequest_Asset{
				Id:       2,
				Required: true,

				AssetOneof: &openrtb.NativeRequest_Asset_Img{
					Img: &openrtb.NativeRequest_Asset_Image{
						Type: openrtb.NativeRequest_MAIN,
					},
				},
			})
		}

		imp := &openrtb.Imp{
			Id:       fmt.Sprintf("cly-%d470%d-%s", pub.ID(), i, <-random.ID),
			Bidfloor: float64(pub.FloorCPM()),
			Secure:   sec,
			Native: &openrtb.Native{
				RequestOneof: &openrtb.Native_RequestNative{},
			},
			Ext: &openrtb.Imp_Ext{
				Mincpc: pub.MinCPC(string(entity.RequestTypeNative)),
			},
		}
		res = append(res, imp)
	}
	return res
}
