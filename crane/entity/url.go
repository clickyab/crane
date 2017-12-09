package entity

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
