package dset

import "time"

// DistributedSet is a set distributed
type DistributedSet interface {
	// Members return all members of this set
	Members() []string
	// Add new item to set
	Add(...string)
	// Key return the master key
	Key() string
	// Save the set with lifetime
	Save(time.Duration)
}

// NewDistributedSet is the distributed set
func NewDistributedSet(name string) DistributedSet {
	return nil
}
