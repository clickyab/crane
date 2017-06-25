package entity

// OS is the os
type OS struct {
	Valid  bool
	Name   string
	Mobile bool
}

func OsFromUA(ua string) OS {
	// Use UA Parser library
	return OS{}
}
