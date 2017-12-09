package cyos

import "strings"

const (
	osMac     int64 = 1
	osUnknown       = 2
	osWindows       = 3
	osLinux         = 4
	osIOS           = 5
	osAndroid       = 6
)

var platforms = map[string]int64{
	"windows":   osWindows,
	"macintosh": osMac,
	"x11":       osLinux,
	"android":   osAndroid,
	"tablet":    osAndroid,
	"iPhone":    osIOS,
	"like Mac":  osIOS,
	"iPod":      osIOS,
	"iPad":      osIOS,
	"linux":     osAndroid,
	"mobile":    osAndroid,
}

// FindOsID try to find os ID base on old id of system
func FindOsID(platform string) int64 {
	if platform == "" {
		return osUnknown
	}
	platform = strings.ToLower(platform)
	p, ok := platforms[platform]
	if ok {
		return p
	}
	return osUnknown
}
