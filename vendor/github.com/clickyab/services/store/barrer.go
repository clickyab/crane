package store

import "time"

// Bearer will use to handle transfer data over requests
type Bearer interface {
	Encode(map[string]string, time.Duration) string
	Decode([]byte, []string) (bool, map[string]string, error)
}
