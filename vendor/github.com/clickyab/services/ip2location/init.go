package ip2location

import (
	"context"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/initializer"
)

type initIP2location struct {
}

func (initIP2location) Initialize(context.Context) {
	assert.Nil(open())
}

func init() {
	initializer.Register(&initIP2location{}, 0)
}
