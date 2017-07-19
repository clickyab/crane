package local

import "github.com/clickyab/services/array"

const (
	linearWidth     int = 800
	linearHeight    int = 440
	nonLinearWidth  int = 728
	nonLinearHeight int = 90
	//LenTypeShort const
	LenTypeShort LenVast = "short"
	//LenTypeDef const
	LenTypeDef LenVast = "def"
	//LenTypeLong const
	LenTypeLong LenVast = "long"
	//TypeVastLinear const
	TypeVastLinear TypeVast = "linear"
	//TypeVastNonLinear const
	TypeVastNonLinear TypeVast = "nonlinear"
)

//LenVast type vast length
type LenVast string

//TypeVast type Vast
type TypeVast string

//IsValidVastLen check valid LenVast
func IsValidVastLen(a string) bool {
	return array.StringInArray(a, string(LenTypeLong), string(LenTypeDef), string(LenTypeShort))
}

// VastSize return size vast
var VastSize = map[TypeVast]map[string]int{
	TypeVastLinear: {
		"width":  linearWidth,
		"height": linearHeight,
	},
	TypeVastNonLinear: {
		"width":  nonLinearWidth,
		"height": nonLinearHeight,
	},
}

var vastConfig = map[string]map[LenVast]map[string][]string{
	"default": {
		"short": {
			"start": {"linear", "11", "00:00:10", "00:00:03"},
			"end":   {"linear", "12", "00:00:10", "00:00:03"},
		},
		"default": {
			"start":    {"linear", "11", "00:00:10", "00:00:03"},
			"00:00:10": {"non-linear", "22", "00:00:06"},
			"end":      {"linear", "13", "00:00:10", "00:00:03"},
		},
		"long": {
			"start":    {"linear", "11", "00:00:10", "00:00:03"},
			"00:00:20": {"non-linear", "22", "00:00:06"},
			"00:01:20": {"non-linear", "23", "00:00:06"},
			"00:03:20": {"non-linear", "24", "00:00:06"},
			"00:05:20": {"non-linear", "25", "00:00:06"},
		},
	},
	"jabeh": {
		"short": {
			"start": {"linear", "11", "00:00:10", "00:00:03"},
			"end":   {"linear", "12", "00:00:10", "00:00:03"},
		},
		"default": {
			"start":    {"linear", "11", "00:00:10", "00:00:03"},
			"00:00:10": {"non-linear", "22", "00:00:06"},
			"end":      {"linear", "13", "00:00:10", "00:00:03"},
		},
		"long": {
			"start":    {"linear", "11", "00:00:10", "00:00:03"},
			"00:00:20": {"non-linear", "22", "00:00:06"},
			"00:01:20": {"non-linear", "23", "00:00:06"},
			"00:03:20": {"non-linear", "24", "00:00:06"},
			"00:05:20": {"non-linear", "25", "00:00:06"},
		},
	},
}

// MakeVastLen return vast len
func MakeVastLen(len LenVast, mod string) map[string][]string {
	if m, found := vastConfig[mod][len]; found {
		return m
	}
	return vastConfig["default"][len]

}
