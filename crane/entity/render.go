package entity

import (
	"net/http"

	"context"

	"clickyab.com/crane/crane/entity"
	"clickyab.com/gad/builder"
)

// Renderer is the app renderer
type Renderer interface {
	// Render render an advertise into a type response, panic if the adType is not supported by advertise
	Render(context.Context, http.ResponseWriter, *builder.Context, entity.Seat, entity.Advertise) error
}
