package entity

const (
	// LinuxOS linux os
	LinuxOS = "linux"
	// AndoidOS Andoid OS
	AndoidOS = "android"
	// WindowsOS Windows OS
	WindowsOS = "windows"
)

// OS is the os
type OS struct {
	Valid  bool
	ID     int64
	Name   string
	Mobile bool
}

func IsMobileOS(name string) bool {
	switch name {
	case AndoidOS:
		return true

	default:
		return false
	}
}

