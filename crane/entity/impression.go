package entity

import (
	"context"

	"clickyab.com/crane/crane/builder"
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
	// Common get common data from context\
	Common() *builder.Common
	// App return th app part of data
	App() *builder.App
	// Data return data of the context
	Data() *builder.Data
	// Data return data of the context
	RTB() *builder.RTB
}
