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
	// Publisher return the publisher that this client is going into system from that
	Publisher() Publisher
	// Slots is the slot for this request
	Seats() []Seat
	// Category returns category obviously
	Category() []Category
	// User return user data
	User() User
	// Tiny means that the logo of clickyab should be shown (true) or not
	Tiny() bool
	// Currency
	Currency() string
	// MultiVideo determine this request can have multiple video
	MultiVideo() bool
	// FloorDiv is floor-cpm divider
	FloorDiv() int64
	// Capping is required otr not
	Capping() bool
}
