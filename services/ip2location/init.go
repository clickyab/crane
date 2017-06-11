package ip2location

import (
	"context"

	"clickyab.com/exchange/services/assert"
	"clickyab.com/exchange/services/initializer"
)

type initIP2location struct {
}

func (initIP2location) Initialize(context.Context) {
	assert.Nil(open())
}

func init() {
	initializer.Register(&initIP2location{}, 0)
}
