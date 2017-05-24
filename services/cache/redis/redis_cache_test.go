package redis

import (
	"testing"

	"encoding/json"
	"io"

	"time"

	"clickyab.com/exchange/services/config"

	"clickyab.com/exchange/services/initializer"

	. "github.com/smartystreets/goconvey/convey"
)

type able struct {
	Key   string
	Test  string
	Value int
}

// Decode try to decode cookie profile into gob
func (cp *able) Decode(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(cp)
}

// Encode try to encode cookie profile from gob
func (cp *able) Encode(i io.Reader) error {
	dnc := json.NewDecoder(i)
	return dnc.Decode(cp)
}

func (cp *able) String() string {
	return cp.Key
}

func TestCacheSystem(t *testing.T) {
	config.Initialize("test", "test", "test")
	defer initializer.Initialize()()

	ch := redisCache{}
	Convey("Test simple cache", t, func() {
		tmp := able{
			Key:   "test_key",
			Test:  "Hi",
			Value: 10,
		}
		So(ch.Do(&tmp, time.Hour), ShouldBeNil)
		var (
			t3 able
		)
		So(ch.Hit(tmp.Key, &t3), ShouldBeNil)
		So(t3, ShouldResemble, tmp)

		So(ch.Hit("NOT VALID KEY", &t3), ShouldNotBeNil)
	})
}
