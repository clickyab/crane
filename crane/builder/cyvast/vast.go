package cyvast

import "strconv"

const (
	short = "short"
	def   = "default"
	long  = "long"
	// VastLinearSize default linear size
	VastLinearSize = 9

	// VastNonLinearSize default non-linear size
	VastNonLinearSize = 6
)

var vastConfig = map[string]map[string][]string{
	"short": {
		"start": {"linear", "11", "00:00:07", "00:00:03"},
		"end":   {"linear", "12", "00:00:10", "00:00:03"},
	},
	"default": {
		"start": {"linear", "11", "00:00:07", "00:00:03"},
		//"00:00:40": {"non-linear", "22", "00:00:06"},
		"end": {"linear", "13", "00:00:10", "00:00:03"},
	},
	"long": {
		"start": {"linear", "11", "00:00:07", "00:00:03"},
		//"00:00:40": {"non-linear", "22", "00:00:06"},
		//"00:01:40": {"non-linear", "23", "00:00:06"},
		//"00:03:40": {"non-linear", "24", "00:00:06"},
		//"00:05:20": {"non-linear", "25", "00:00:06"},
	},
}

// MakeVastLen return vast len
func MakeVastLen(len string, first bool, mid bool, end bool) (string, map[string][]string) {
	//apply vast customization
	preRes := vastConfig
	for i := range preRes {
		if !first {
			delete(preRes[i], "start")
		}
		if !end {
			delete(preRes[i], "end")
		}
		if !mid {
			for j := range preRes[i] {
				if j != "start" && j != "end" {
					delete(preRes[i], j)
				}
			}
		}
	}
	if m, found := preRes[len]; found {
		return def, m
	}
	if m, err := strconv.ParseInt(len, 10, 64); err == nil {
		if m < 30 {
			return short, preRes[short]
		} else if m < 90 {
			return def, preRes[def]
		} else {
			return long, preRes[long]
		}
	}
	return def, preRes[def]
}
