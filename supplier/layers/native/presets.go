package native

import (
	"fmt"
	"net/http"

	"github.com/golang/protobuf/jsonpb"

	"clickyab.com/crane/demand/entity"
	openrtb "clickyab.com/crane/openrtb/v2.5"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/random"
)

// Version for default version of native request
var Version = 1

func getImps(r *http.Request, count int, pub entity.Publisher, image bool) []*openrtb.Imp {
	var (
		res []*openrtb.Imp
		sec = framework.Scheme(r) == "https"
	)
	for i := 1; i <= count; i++ {
		// make request
		req := &openrtb.NativeRequest{
			Assets: []*openrtb.NativeRequest_Asset{
				{
					Id:       1,
					Required: 1,
					AssetOneof: &openrtb.NativeRequest_Asset_Title_{
						Title: &openrtb.NativeRequest_Asset_Title{
							Len: int32(nativeMaxTitleLen.Int()),
						},
					},
				},
			},
		}
		if image {
			req.Assets = append(req.Assets, &openrtb.NativeRequest_Asset{
				Id:       2,
				Required: 1,
				AssetOneof: &openrtb.NativeRequest_Asset_Img{
					Img: &openrtb.NativeRequest_Asset_Image{
						Type: openrtb.NativeRequest_MAIN,
					},
				},
			})
		}

		if Version == 0 {

			jm := jsonpb.Marshaler{}
			s, err := jm.MarshalToString(req)
			if err != nil {
				fmt.Println(err)
			}

			imp := &openrtb.Imp{
				Id:       fmt.Sprintf("cly-%d470%d-%s", pub.ID(), i, <-random.ID),
				Bidfloor: float64(pub.FloorCPM()),
				Secure: func() int32 {
					if sec {
						return 1
					}
					return 0
				}(),
				Native: &openrtb.Native{
					RequestOneof: &openrtb.Native_Request{
						Request: s,
					},
				},
				Ext: &openrtb.Imp_Ext{
					Mincpc: pub.MinCPC(string(entity.RequestTypeNative)),
				},
			}

			res = append(res, imp)
			continue
		}

		imp := &openrtb.Imp{
			Id:       fmt.Sprintf("cly-%d470%d-%s", pub.ID(), i, <-random.ID),
			Bidfloor: float64(pub.FloorCPM()),
			Secure: func() int32 {
				if sec {
					return 1
				}
				return 0
			}(),
			Native: &openrtb.Native{
				RequestOneof: &openrtb.Native_RequestNative{
					RequestNative: req,
				},
			},
			Ext: &openrtb.Imp_Ext{
				Mincpc: pub.MinCPC(string(entity.RequestTypeNative)),
			},
		}

		res = append(res, imp)

	}
	return res
}
