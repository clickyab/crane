package cache_test

import (
	"testing"

	. "clickyab.com/exchange/services/cache"
	"clickyab.com/exchange/services/cache/mock"

	"encoding/json"
	"io"

	"time"

	"errors"

	. "github.com/smartystreets/goconvey/convey"
)

type able struct {
	Key   string
	Test  string
	Value int
}

type notable struct {
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
	ch := mock.NewCacheMock()
	Register(ch)
	Convey("Test simple cache", t, func() {
		tmp := able{
			Key:   "test_key",
			Test:  "Hi",
			Value: 10,
		}
		So(Do(&tmp, time.Hour, errors.New("example")), ShouldNotBeNil)
		So(Do(&tmp, time.Hour, nil), ShouldBeNil)
		has, v := mock.GetData(ch, tmp.Key)
		So(has, ShouldBeTrue)
		var (
			t2 able
			t3 able
		)
		So(json.Unmarshal(v, &t2), ShouldBeNil)
		So(t2, ShouldResemble, tmp)

		So(Hit(tmp.Key, &t3), ShouldBeNil)
		So(t3, ShouldResemble, tmp)

		So(Hit("NOT VALID KEY", &t3), ShouldNotBeNil)
	})
	ch = mock.NewCacheMock()
	Register(ch)
	Convey("Test cache wrapper", t, func() {
		tmp := notable{
			Key:   "test_key",
			Test:  "Hi",
			Value: 10,
		}
		wrap := CreateWrapper(tmp.Key, &tmp)
		So(Do(wrap, time.Hour, errors.New("example")), ShouldNotBeNil)
		So(Do(wrap, time.Hour, nil), ShouldBeNil)
		has, _ := mock.GetData(ch, tmp.Key)
		So(has, ShouldBeTrue)
		var (
			t3 notable
		)
		wrap2 := CreateWrapper(tmp.Key, &t3)
		So(Hit(tmp.Key, wrap2), ShouldBeNil)
		So(t3, ShouldResemble, tmp)
		So(wrap2.Entity(), ShouldResemble, &tmp)

		So(Hit("NOT VALID KEY", CreateWrapper(tmp.Key, &t3)), ShouldNotBeNil)
	})
}
