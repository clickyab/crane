package entity

import (
	"io"

	"context"
)

// Renderer is the app renderer
type Renderer interface {
	// Render render an advertise into a type response, panic if the adType is not supported by advertise
	Render(context.Context, io.Writer, Context, Seat) error
}
