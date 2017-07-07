package entity

// ClickProvider is the provider for urls used in the system
type ClickProvider interface {
	// ClickURL return a click link for app
	ClickURL(Slot, Impression) string
}

// ClickStatus determines status of click
type ClickStatus int

const (

	// SuspSuccessful status
	SuspSuccessful = ClickStatus(iota + 1e2)

	// SuspFastClick status
	SuspFastClick

	// SuspNoAdFound status
	SuspNoAdFound

	// SuspIPMismatch status
	SuspIPMismatch

	// SuspUAMismatch status
	SuspUAMismatch
)
