package redis

import (
	"testing"

	"encoding/json"
	"io"

	"time"

	"github.com/clickyab/services/config"

	"github.com/clickyab/services/initializer"

	. "github.com/smartystreets/goconvey/convey"
)

type able struct {
	Key   string
	Test  string
	Value int
}

// Encode try to decode cookie profile into gob
func (cp *able) Encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(cp)
}

// Decode try to encode cookie profile from gob
func (cp *able) Decode(i io.Reader) error {
	dnc := json.NewDecoder(i)
	return dnc.Decode(cp)
}

func (cp *able) String() string {
	return cp.Key
}

func TestCacheSystem(t *testing.T) {
	config.Initialize("test", "test", "test")
	defer initializer.Initialize()()

	ch := cache{}
	Convey("Test simple cache", t, func() {
		tmp := able{
			Key:   "test_key",
			Test:  "Hi",
			Value: 10,
		}
		So(ch.Do(tmp.String(), &tmp, time.Hour), ShouldBeNil)
		var (
			t3 able
		)
		So(ch.Hit(tmp.Key, &t3), ShouldBeNil)
		So(t3, ShouldResemble, tmp)

		So(ch.Hit("NOT VALID KEY", &t3), ShouldNotBeNil)
	})
}
