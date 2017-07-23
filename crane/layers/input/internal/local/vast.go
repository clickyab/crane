package local

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	os2 "os"

	"github.com/Sirupsen/logrus"
	"github.com/clickyab/services/array"
	"github.com/clickyab/services/config"
)

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

var vastConfig map[string]map[LenVast]map[string][]string

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

// MakeVastLen return vast len
func MakeVastLen(len LenVast, mod string) map[string][]string {
	vastFile := config.RegisterString("crane.vast.file", "/bin/build/vast_config.json", "read config vast from file").String()
	raw, err := ioutil.ReadFile(fmt.Sprintf("%s%s", os2.Getenv("PWD"), vastFile))
	if err != nil {
		logrus.Debug("can't find config vast file")
	}
	err = json.Unmarshal(raw, &vastConfig)
	if err != nil {
		logrus.Debug("structure json vast file have problem")
	}
	if m, found := vastConfig[mod][len]; found {
		return m
	}
	return vastConfig["default"][len]

}
