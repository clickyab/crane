package restful

import (
	"errors"

	"encoding/json"

	"clickyab.com/crane/crane/entity"
	"github.com/clickyab/services/assert"
)

// New validate request and return new vast impression
func New(r entity.Request) (entity.Context, error) {
	m := &impression{}

	v, ok := r.Attributes()["body"]
	if !ok {
		return nil, errors.New("Attributes does not contains json body")
	}
	e := json.Unmarshal([]byte(v), m)
	assert.Nil(e)

	// TODO: Validate restful input data
	// return m, nil
	return nil, errors.New("Validate restful input data")

}
