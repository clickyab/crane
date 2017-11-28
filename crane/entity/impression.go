package entity

import (
	"context"
)

// ImpressionAttributes is the imp attr key
type ImpressionAttributes string

// ImpressionLayer interface to handle the input layer
type ImpressionLayer interface {
	New(context.Context, Request) (Context, error)
}

// Context is the single impression object
type Context interface {
	// Request data comes from request for every user
	// like ip,user agent,client id,...
	Request
	// TrackID return the random id of this imp object
	TrackID() string
	// Publisher return the publisher that this client is going into system from that
	Publisher() Publisher
	// Slots is the slot for this request
	Slots() []Slot
	// Category returns category obviously
	Category() []Category
}
