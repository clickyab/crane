package vast

import (
	"errors"
	"fmt"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
	"clickyab.com/crane/crane/models/query"
	"github.com/clickyab/services/random"
)

var (
	// ErrorPublisherNotFound error when query can't find publisher
	ErrorPublisherNotFound = errors.New("publisher with the specified domain not found")
	// ErrorPublisherSupplierEmpty error when publisher or supplier is wmpty
	ErrorPublisherSupplierEmpty = errors.New("publisher or supplier empty")
	// ErrorLenVast error len vast
	ErrorLenVast = errors.New("len for vast can't be empty or is not standard")
)

// New validate request and return a vast impression
func New(r entity.Request) (entity.Impression, error) {
	//fetch website by domain and (clickyab) supplier
	lenVast := r.Attributes()["l"]
	if lenVast == "" || !local.IsValidVastLen(lenVast) {
		return nil, ErrorLenVast
	}
	if r.Attributes()["d"] == "" || r.Attributes()["supplier"] == "" {
		return nil, ErrorPublisherSupplierEmpty
	}
	q := query.Publisher()
	pub, err := q.ByPlatform(r.Attributes()["d"], entity.WebPlatform, r.Attributes()["supplier"])
	if err != nil {
		return nil, ErrorPublisherNotFound
	}

	res := &impression{}
	under := pub.UnderFloor()

	res.pub = &local.Publisher{
		FName:         pub.Name(),
		FSupplier:     pub.Supplier(),
		FFloorCPM:     pub.FloorCPM(),
		FSoftFloorCPM: pub.SoftFloorCPM(),
		FUnderFloor:   &under,
	}

	res.trackID = <-random.ID
	mod := r.Attributes()["mod"]
	vastConf := local.MakeVastLen(local.LenVast(lenVast), mod)
	sResult := make([]entity.Slot, len(vastConf))
	attr := make(map[string]interface{})
	for i := range vastConf {
		attr["offset"] = i
		attr["breakType"] = vastConf[i][0]
		attr["duration"] = vastConf[i][2]
		if vastConf[i][0] == string(local.TypeVastLinear) {
			attr["repeat"] = vastConf[i][3]
		}
		attr["mod"] = mod
		attr["len"] = lenVast
		tVast, ok := attr["breakType"].(local.TypeVast)
		if !ok {
			continue
		}
		sResult = append(sResult, local.ExtractSlot(pub.Supplier(), pub.Name(), entity.VastPlatform, local.VastSize[tVast]["width"], local.VastSize[tVast]["height"], fmt.Sprintf("%d", i), attr))
	}
	res.slots = sResult
	// TODO implement later
	res.categories = make([]entity.Category, 0)
	return res, nil

}

func (i *impression) Attributes() map[string]string {
	return i.attr
}
