package entity

import (
	"io"
)

// Renderer is the app renderer
type Renderer interface {
	// Render render an advertise into a type response, panic if the adType is not supported by advertise
	Render(Advertise, AdType, io.Writer, URLProvider) error
}
