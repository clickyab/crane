package native

import (
	"errors"
	"fmt"
	"strconv"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
	"clickyab.com/crane/crane/models/query"
	"github.com/clickyab/services/random"
)

// New validate request and return a native impression
func New(r entity.Request) (entity.Impression, error) {
	//fetch website by domain and (clickyab) supplier
	q := query.Publisher()
	if r.Attributes()["d"] == "" || r.Attributes()["supplier"] == "" {
		return nil, errors.New("domain or supplier empty")
	}
	pub, err := q.ByPlatform(r.Attributes()["d"], entity.WebPlatform, r.Attributes()["supplier"])
	if err != nil {
		return nil, errors.New("publisher with the specified domain not found")
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

	intCount, err := strconv.Atoi(r.Attributes()["count"])
	if err != nil {
		return nil, errors.New("count not valid")
	}
	sResult := make([]entity.Slot, intCount)
	attr := make(map[string]interface{})
	for i := 1; i <= intCount; i++ {
		//fill slot
		sResult = append(sResult, local.ExtractSlot(pub.Supplier(), pub.Name(), entity.NativePlatform, 0, 0, fmt.Sprintf("%d", i), attr))
	}
	res.slots = sResult
	// TODO implement later
	res.categories = make([]entity.Category, 0)
	return res, nil

}

func (i *impression) Attributes() map[string]string {
	return i.attr
}
