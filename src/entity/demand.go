package entity

import (
	"context"
	"net/http"
)

// Demand is the interface to handle ad in system base on impression
type Demand interface {
	// Name return the name of this demand
	Name() string
	// Provide is the function to handle the request, provider should response
	// to this function and return all eligible ads
	// a very important note about providers :
	// they must return as soon as possible, waiting for result must be done
	// in separate co-routine. just create a channel, run a co-routine, and return.
	Provide(context.Context, Impression, chan map[string]Advertise)
	// Win return the win response to the demand. it happen only if the request is the winner
	Win(context.Context, string)
	// Status is called for getting the statistics of this Demand
	Status(context.Context, http.ResponseWriter, *http.Request)
}
