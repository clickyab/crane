package mock

import (
	"testing"

	"encoding/json"
	"io"

	"time"

	"errors"

	"github.com/clickyab/services/kv"
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
	ch := NewCacheMock()
	kv.Register(nil, nil, nil, nil, nil, ch, nil, nil)
	Convey("Test simple cacheFactory", t, func() {
		tmp := able{
			Key:   "test_key",
			Test:  "Hi",
			Value: 10,
		}
		So(kv.Do(tmp.Key, &tmp, time.Hour, errors.New("example")), ShouldNotBeNil)
		So(kv.Do(tmp.Key, &tmp, time.Hour, nil), ShouldBeNil)
		has, v := GetData(ch, tmp.Key)
		So(has, ShouldBeTrue)
		var (
			t2 able
			t3 able
		)
		So(json.Unmarshal(v, &t2), ShouldBeNil)
		So(t2, ShouldResemble, tmp)

		So(kv.Hit(tmp.Key, &t3), ShouldBeNil)
		So(t3, ShouldResemble, tmp)

		So(kv.Hit("NOT VALID KEY", &t3), ShouldNotBeNil)
	})
	ch = NewCacheMock()
	kv.Register(nil, nil, nil, nil, nil, ch, nil, nil)
	Convey("Test cacheFactory wrapper", t, func() {
		tmp := notable{
			Key:   "test_key",
			Test:  "Hi",
			Value: 10,
		}
		wrap := kv.CreateWrapper(&tmp)
		So(kv.Do(tmp.Key, wrap, time.Hour, errors.New("example")), ShouldNotBeNil)
		So(kv.Do(tmp.Key, wrap, time.Hour, nil), ShouldBeNil)
		has, _ := GetData(ch, tmp.Key)
		So(has, ShouldBeTrue)
		var (
			t3 notable
		)
		wrap2 := kv.CreateWrapper(&t3)
		So(kv.Hit(tmp.Key, wrap2), ShouldBeNil)
		So(t3, ShouldResemble, tmp)
		So(wrap2.Entity(), ShouldResemble, &tmp)

		So(kv.Hit("NOT VALID KEY", kv.CreateWrapper(&t3)), ShouldNotBeNil)
	})
}
