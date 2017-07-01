package entity

import (
	"github.com/mssola/user_agent"
)

// OS is the os
type OS struct {
	Valid  bool
	Name   string
	Mobile bool
}

// OsFromUA return os
func OsFromUA(ua string) OS {
	client := user_agent.New(ua)
	return OS{
		Valid:  true,
		Name:   client.OS(),
		Mobile: client.Mobile(),
	}
}
