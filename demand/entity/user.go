package entity

// User in system
type User interface {
	// ID return user id
	ID() string

	// CyuId return unique id for page
	CyuId() string
}
