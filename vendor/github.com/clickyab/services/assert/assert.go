package assert

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func doPanic(err string, tag string, params ...interface{}) {
	f := logrus.Fields{}
	for i := range params {
		f[fmt.Sprintf("param%d", i)] = params[i]
	}
	f["tag"] = tag
	logrus.WithFields(f).Panic(err)
}

// Nil panic if the test is not nil
func Nil(test interface{}, params ...interface{}) {
	if test != nil {
		tag := "panic.nil"
		if e, ok := test.(error); ok {
			doPanic(e.Error(), tag, params...)
			return
		}
		doPanic("must be nil but is not", tag, params...)
	}
}

// NotNil panic if the test is nil
func NotNil(test interface{}, params ...interface{}) {
	if test == nil {
		doPanic("must not be nil, but it is", "panic.notnil", params...)
	}
}

// True check if the value is true and panic if its not
func True(test bool, params ...interface{}) {
	if !test {
		doPanic("must be true, but its not", "panic.true", params...)
	}
}

// False check if the test is false
func False(test bool, params ...interface{}) {
	if test {
		doPanic("must be false, but its not", "panic.false", params...)
	}

}

// Empty check if the string is empty and panic if not
func Empty(test string, params ...interface{}) {
	if test != "" {
		doPanic("must be empty, but its not", "panic.empty", params...)
	}
}
