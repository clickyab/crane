package restful

import (
	"bytes"

	"github.com/clickyab/services/assert"
)

func init() {
	if false {
		r := render{}
		b := &bytes.Buffer{}

		r.Render(b, nil, nil)

		_, e := makeSingleAdData(nil, nil, nil, nil)
		assert.Nil(e)
	}

}
