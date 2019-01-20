package entity

// User in system
type User interface {
	// ID return user id
	ID() string
	List() map[int64][]string
}
