package store

import "time"

// Bearer will use to handle transfer data over requests
type Bearer interface {
	// Encode a map into a string, with time key
	Encode(map[string]string, time.Duration) string
	// Decode a map and return the decoded map, also check for expiration, this
	// suppose to work with encoded value with the same interface not any encoded value
	Decode([]byte, ...string) (bool, map[string]string, error)
}
