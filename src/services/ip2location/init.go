package ip2location

import (
	"context"
	"services/initializer"
)

type initIP2location struct {
}

func (initIP2location) Initialize(context.Context) {
	open()
}

func init() {
	initializer.Register(&initIP2location{}, 0)
}
