package restful

import (
	"errors"

	"encoding/json"

	"fmt"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/crane/crane/layers/input/internal/local"
	"clickyab.com/crane/crane/models/query"
	"github.com/clickyab/services/assert"
)

var (
	// ErrorPublisherNotFound error when query can't find publisher
	ErrorPublisherNotFound = errors.New("publisher with the specified domain not found")
	// ErrorPublisherSupplierEmpty error when publisher or supplier is wmpty
	ErrorPublisherSupplierEmpty = errors.New("publisher or supplier empty")
	// ErrorAttrBodyEmpty error when attr body empty
	ErrorAttrBodyEmpty = errors.New("Attributes does not contains json body")
)

// New validate request and return new rest impression
func New(r entity.Request) (entity.Impression, error) {
	m := &rawImpression{}
	v, ok := r.Attributes()["body"]
	if !ok {
		return nil, ErrorAttrBodyEmpty
	}
	e := json.Unmarshal([]byte(v), m)
	assert.Nil(e)
	// check publisher
	if m.Publisher().Name() == "" || m.Publisher().Supplier() == "" {
		return nil, ErrorPublisherSupplierEmpty
	}
	q := query.Publisher()
	publisher, err := q.ByPlatform(m.Publisher().Name(), entity.WebPlatform, m.Publisher().Supplier())
	if err != nil {
		return nil, ErrorPublisherNotFound
	}
	under := publisher.UnderFloor()
	m.publisher = &local.Publisher{
		FName:         publisher.Name(),
		FSupplier:     publisher.Supplier(),
		FFloorCPM:     publisher.FloorCPM(),
		FSoftFloorCPM: publisher.SoftFloorCPM(),
		FUnderFloor:   &under,
	}
	for i, x := range m.FSlots {
		m.FSlots[i].FID = fmt.Sprintf("%s/%s/%s/%dx%d/%s", m.Publisher().Supplier(), m.Publisher().Name(), entity.WebPlatform, x.Width(), x.Height(), x.TrackID())
	}

	return m, nil

}
