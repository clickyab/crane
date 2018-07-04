package entities

import "time"

// StaticSeat static seat interface
type StaticSeat interface {
	ID() int64
	Publisher() string
	Supplier() string
	Type() string
	Position() string
	From() time.Time
	To() time.Time
	Fix() bool
	RTBMarkup() string
	Chance() int
}
